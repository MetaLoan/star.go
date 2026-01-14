import { useEffect, useState, useCallback, useRef } from "react"
import { X, ChevronLeft, ChevronRight, Plus, Camera, RefreshCw, Scan, CheckCircle2 } from "lucide-react"
import { InkRevealText } from "./InkRevealText"
import { SolarSystem } from "./SolarSystem"
import { OracleInteractiveEye } from "./OracleInteractiveEye"
import { cn } from "../../utils/cn"

type Phase = "choice" | "creating" | "editing" | "resonating" | "doorOpening" | "result" | "exiting"
type SelectionStep = "first" | "second"

interface UserOrb {
  id: string
  name: string
  gender: string
  birthday: string
  soul: string
  blood: string
}

interface ConnectUnifiedProps {
  isVisible: boolean
  onAccept: () => void
  onReject: () => void
  onComplete: (result: any) => void
  onExit: () => void
}

const mockUsers: UserOrb[] = [
  { id: "1", name: "Leo", gender: "男", birthday: "1990-08-15", soul: "SOUL-A1B2C3D4-XY12", blood: "O" },
  { id: "2", name: "Luna", gender: "女", birthday: "1995-03-22", soul: "", blood: "A" },
  { id: "3", name: "Nova", gender: "", birthday: "1988-11-11", soul: "", blood: "" },
  { id: "4", name: "Aria", gender: "", birthday: "", soul: "", blood: "" },
]

const ringPositions = Array.from({ length: 5 }, (_, index) => {
  const angleDeg = -90 + index * (360 / 5)
  const angleRad = (angleDeg * Math.PI) / 180
  const radius = 32
  const x = 50 + radius * Math.cos(angleRad)
  const y = 50 + radius * Math.sin(angleRad)
  return { top: `${y}%`, left: `${x}%`, x, y }
})

const starOrder = [0, 2, 4, 1, 3]
const starPoints = starOrder.map((index) => `${ringPositions[index].x},${ringPositions[index].y}`).join(" ")

const ZodiacIcons = {
  Aries: (props: any) => (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" {...props}>
      <path d="M12 21c0-4.5 0-9 0-9m0 0c0-3 2-5 5-5s5 2 5 5m-10 0c0-3-2-5-5-5s-5 2-5 5" />
    </svg>
  ),
  Taurus: (props: any) => (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" {...props}>
      <circle cx="12" cy="14" r="6" />
      <path d="M5 4c1 4 3 6 7 6s6-2 7-6" />
    </svg>
  ),
  Gemini: (props: any) => (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" {...props}>
      <path d="M7 4a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2m10-16a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2M4 4h16M4 20h16" />
    </svg>
  ),
  Cancer: (props: any) => (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" {...props}>
      <circle cx="16" cy="7" r="3" />
      <circle cx="8" cy="17" r="3" />
      <path d="M19 10c-3 0-6 2-7 5m-7-3c3 0 6-2 7-5" />
    </svg>
  ),
  Leo: (props: any) => (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" {...props}>
      <circle cx="6" cy="17" r="3" />
      <path d="M8.5 15c1-3 4-5 7-3s3 5 0 8m0-13c2 0 4 2 4 4s-2 4-4 4" />
    </svg>
  ),
  Virgo: (props: any) => (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" {...props}>
      <path d="M5 4v12a3 3 0 0 0 6 0V4m0 12a3 3 0 0 0 6 0V4m0 12c0 2 1 4 3 4" />
      <path d="M11 4h1m-7 0h1" />
    </svg>
  ),
  Libra: (props: any) => (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" {...props}>
      <path d="M5 20h14M5 16h14m-7-4c-3 0-5-2-5-5h10c0 3-2 5-5 5z" />
    </svg>
  ),
  Scorpio: (props: any) => (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" {...props}>
      <path d="M5 4v10a2 2 0 0 0 4 0V4m0 10a2 2 0 0 0 4 0V4m0 10a2 2 0 0 0 4 0m0 0c0 2 1 3 3 3" />
      <path d="M17 17l2 2m-2 0l2-2" />
    </svg>
  ),
  Sagittarius: (props: any) => (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" {...props}>
      <path d="M4 20L20 4m0 0h-7m7 0v7M7 13l4 4" />
    </svg>
  ),
  Capricorn: (props: any) => (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" {...props}>
      <path d="M4 5l4 14 4-8c1-2 4-2 4 0s-2 6-2 8 2 2 4 0" />
    </svg>
  ),
  Aquarius: (props: any) => (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" {...props}>
      <path d="M4 10l3-3 4 4 3-3 4 4 2-2M4 17l3-3 4 4 3-3 4 4 2-2" />
    </svg>
  ),
  Pisces: (props: any) => (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" {...props}>
      <path d="M4 12c5 0 8-3 8-8m0 16c0-5 3-8 8-8M12 4v16m-8 0c0-5 3-8 8-8m0 0c5 0 8 3 8 8" />
    </svg>
  ),
}

const zodiacSigns = [
  { name: "Capricorn", icon: ZodiacIcons.Capricorn, start: [1, 1], end: [1, 19] },
  { name: "Aquarius", icon: ZodiacIcons.Aquarius, start: [1, 20], end: [2, 18] },
  { name: "Pisces", icon: ZodiacIcons.Pisces, start: [2, 19], end: [3, 20] },
  { name: "Aries", icon: ZodiacIcons.Aries, start: [3, 21], end: [4, 19] },
  { name: "Taurus", icon: ZodiacIcons.Taurus, start: [4, 20], end: [5, 20] },
  { name: "Gemini", icon: ZodiacIcons.Gemini, start: [5, 21], end: [6, 20] },
  { name: "Cancer", icon: ZodiacIcons.Cancer, start: [6, 21], end: [7, 22] },
  { name: "Leo", icon: ZodiacIcons.Leo, start: [7, 23], end: [8, 22] },
  { name: "Virgo", icon: ZodiacIcons.Virgo, start: [8, 23], end: [9, 22] },
  { name: "Libra", icon: ZodiacIcons.Libra, start: [9, 23], end: [10, 22] },
  { name: "Scorpio", icon: ZodiacIcons.Scorpio, start: [10, 23], end: [11, 21] },
  { name: "Sagittarius", icon: ZodiacIcons.Sagittarius, start: [11, 22], end: [12, 21] },
  { name: "Capricorn", icon: ZodiacIcons.Capricorn, start: [12, 22], end: [12, 31] },
]

function getZodiacFromBirthday(birthday: string): { name: string; icon: any } | null {
  if (!birthday) return null
  const date = new Date(birthday)
  if (isNaN(date.getTime())) return null
  const month = date.getMonth() + 1
  const day = date.getDate()
  for (const sign of zodiacSigns) {
    const [startMonth, startDay] = sign.start
    const [endMonth, endDay] = sign.end
    if (startMonth === endMonth) {
      if (month === startMonth && day >= startDay && day <= endDay) return { name: sign.name, icon: sign.icon }
    } else {
      if ((month === startMonth && day >= startDay) || (month === endMonth && day <= endDay)) return { name: sign.name, icon: sign.icon }
    }
  }
  return null
}

const MaterialIcon = ({ name, className = "" }: { name: string; className?: string }) => (
  <span 
    className={cn("material-symbols-outlined transition-opacity duration-300", className)}
    style={{ 
      fontSize: '26px',
      color: 'inherit',
      fontVariationSettings: "'FILL' 0, 'wght' 200, 'GRAD' 0, 'opsz' 24"
    }}
  >
    {name}
  </span>
)

const nodeConfigs = [
  { label: "姓名", field: "name", icon: "person" },
  { label: "性别", field: "gender", icon: "person_2" },
  { label: "生日", field: "birthday", icon: "calendar_today" },
  { label: "灵魂", field: "soul", icon: "visibility_off" },
  { label: "血型", field: "blood", icon: "water_drop" },
]

function EyeScanner({ onCapture }: { onCapture: (code: string, image: string) => void }) {
  const [isStarted, setIsStarted] = useState(false)
  const [hasPermission, setHasPermission] = useState<boolean | null>(null)
  const [loadingStatus, setLoadingStatus] = useState<string>("IDLE")
  const [eyePositions, setEyePositions] = useState<{ left: { x: number, y: number }, right: { x: number, y: number } } | null>(null)
  const [capturePhase, setCapturePhase] = useState<number>(-1)
  const capturePhaseRef = useRef<number>(-1)
  const [capturedImage, setCapturedImage] = useState<string | null>(null)
  const [capturedEyePositions] = useState<{ left: { x: number, y: number }, right: { x: number, y: number } } | null>(null)
  const [soulCode, setSoulCode] = useState<string | null>(null)
  
  const videoRef = useRef<HTMLVideoElement>(null)
  const canvasRef = useRef<HTMLCanvasElement>(null)
  const streamRef = useRef<MediaStream | null>(null)
  const faceMeshRef = useRef<any>(null)
  const requestRef = useRef<number | null>(null)
  const isActiveRef = useRef(false)
  const eyeDetectedRef = useRef(false)

  const ritualPhrases = [
    "正在透过灵魂之门凝视...",
    "解读虹膜中的远古图案...",
    "解码你本质的星象印记...",
    "感应你存在的以太频率...",
    "捕捉你内在宇宙的无限深度..."
  ]

  const waitingMessage = "等待你的凝视... 请让我看到你的双眼正视镜头。"
  const LEFT_IRIS_CENTER = 468
  const RIGHT_IRIS_CENTER = 473

  useEffect(() => {
    capturePhaseRef.current = capturePhase
  }, [capturePhase])

  useEffect(() => {
    if (!eyeDetectedRef.current || capturePhase < 0 || capturePhase >= ritualPhrases.length || capturedImage) {
      return
    }
    const timer = setTimeout(() => {
      if (eyeDetectedRef.current) {
        const nextPhase = capturePhase + 1
        if (nextPhase < ritualPhrases.length) {
          setCapturePhase(nextPhase)
        } else {
          performCapture()
        }
      }
    }, 2000)
    return () => clearTimeout(timer)
  }, [capturePhase, capturedImage])

  const waitForLibs = async (maxRetries = 30): Promise<boolean> => {
    for (let i = 0; i < maxRetries; i++) {
      const fm = (window as any).FaceMesh
      if (fm) return true
      await new Promise(r => setTimeout(r, 300))
      setLoadingStatus(`初始化视觉核心 (${i + 1}/${maxRetries})`)
    }
    return false
  }

  const loadMediaPipeScript = () => {
    return new Promise<void>((resolve, reject) => {
      if ((window as any).FaceMesh) {
        resolve()
        return
      }
      const cameraScript = document.createElement('script')
      cameraScript.src = 'https://cdn.jsdelivr.net/npm/@mediapipe/camera_utils/camera_utils.js'
      cameraScript.crossOrigin = 'anonymous'
      document.head.appendChild(cameraScript)
      const drawingScript = document.createElement('script')
      drawingScript.src = 'https://cdn.jsdelivr.net/npm/@mediapipe/drawing_utils/drawing_utils.js'
      drawingScript.crossOrigin = 'anonymous'
      document.head.appendChild(drawingScript)
      const script = document.createElement('script')
      script.src = 'https://cdn.jsdelivr.net/npm/@mediapipe/face_mesh/face_mesh.js'
      script.crossOrigin = 'anonymous'
      script.onload = () => resolve()
      script.onerror = () => reject(new Error('Failed to load MediaPipe'))
      document.head.appendChild(script)
    })
  }

  const performCapture = () => {
    if (!videoRef.current) return
    const video = videoRef.current
    const vw = video.videoWidth
    const vh = video.videoHeight
    let minX = 0, maxX = vw, minY = 0, maxY = vh
    if (eyePositions && eyePositions.left && eyePositions.right) {
      const leftEyeX = (eyePositions.left.x / 100) * vw
      const leftEyeY = (eyePositions.left.y / 100) * vh
      const rightEyeX = (eyePositions.right.x / 100) * vw
      const rightEyeY = (eyePositions.right.y / 100) * vh
      const paddingX = vw * 0.15
      const paddingY = vh * 0.08
      minX = Math.max(0, Math.min(leftEyeX, rightEyeX) - paddingX)
      maxX = Math.min(vw, Math.max(leftEyeX, rightEyeX) + paddingX)
      minY = Math.max(0, Math.min(leftEyeY, rightEyeY) - paddingY)
      maxY = Math.min(vh, Math.max(leftEyeY, rightEyeY) + paddingY)
    } else {
      const centerX = vw / 2
      const centerY = vh / 2
      const cropW = vw * 0.6
      const cropH = vh * 0.3
      minX = centerX - cropW / 2
      maxX = centerX + cropW / 2
      minY = centerY - cropH / 2
      maxY = centerY + cropH / 2
    }
    const cropWidth = Math.round(maxX - minX)
    const cropHeight = Math.round(maxY - minY)
    const tempCanvas = document.createElement('canvas')
    tempCanvas.width = cropWidth
    tempCanvas.height = cropHeight
    const tempCtx = tempCanvas.getContext('2d')
    if (tempCtx) {
      tempCtx.translate(cropWidth, 0)
      tempCtx.scale(-1, 1)
      tempCtx.drawImage(video, Math.round(minX), Math.round(minY), cropWidth, cropHeight, 0, 0, cropWidth, cropHeight)
      const croppedImageData = tempCanvas.toDataURL('image/jpeg', 0.95)
      setCapturedImage(croppedImageData)
      const hash = croppedImageData.length.toString(16).toUpperCase().slice(-8)
      const code = `SOUL-${hash}-${Date.now().toString(36).toUpperCase().slice(-4)}`
      setSoulCode(code)
      onCapture(code, croppedImageData)
      stopCamera()
    }
  }

  const startCamera = async () => {
    setIsStarted(true)
    setLoadingStatus("加载视觉组件...")
    try {
      await loadMediaPipeScript()
      const libsReady = await waitForLibs()
      if (!libsReady) throw new Error("MEDIAPIPE_LIBS_NOT_FOUND")
      setLoadingStatus("连接摄像头...")
      const stream = await navigator.mediaDevices.getUserMedia({ video: { facingMode: "user", width: { ideal: 640 }, height: { ideal: 480 } }, audio: false })
      streamRef.current = stream
      if (videoRef.current) {
        videoRef.current.srcObject = stream
        await videoRef.current.play()
      }
      setHasPermission(true)
      setLoadingStatus("初始化神经网络...")
      const FaceMesh = (window as any).FaceMesh
      const faceMesh = new FaceMesh({ locateFile: (file: string) => `https://cdn.jsdelivr.net/npm/@mediapipe/face_mesh/${file}` })
      faceMesh.setOptions({ maxNumFaces: 1, refineLandmarks: true, minDetectionConfidence: 0.5, minTrackingConfidence: 0.5 })
      faceMesh.onResults((results: any) => {
        if (results.multiFaceLandmarks && results.multiFaceLandmarks.length > 0) {
          const landmarks = results.multiFaceLandmarks[0]
          const leftEye = landmarks[LEFT_IRIS_CENTER]
          const rightEye = landmarks[RIGHT_IRIS_CENTER]
          if (leftEye && rightEye) {
            setEyePositions({ left: { x: leftEye.x * 100, y: leftEye.y * 100 }, right: { x: rightEye.x * 100, y: rightEye.y * 100 } })
            if (!eyeDetectedRef.current && capturePhase < 0 && !capturedImage) {
              eyeDetectedRef.current = true
              setCapturePhase(0)
            } else {
              eyeDetectedRef.current = true
            }
          } else {
            handleEyesLost()
          }
        } else {
          handleEyesLost()
        }
        renderFrame(results)
        if (loadingStatus !== "READY") setLoadingStatus("READY")
      })
      faceMeshRef.current = faceMesh
      isActiveRef.current = true
      const processFrame = async () => {
        if (!isActiveRef.current || !videoRef.current || videoRef.current.paused) {
          requestRef.current = requestAnimationFrame(processFrame)
          return
        }
        try { await faceMeshRef.current.send({ image: videoRef.current }) } catch (e) { console.error("Frame processing error:", e) }
        requestRef.current = requestAnimationFrame(processFrame)
      }
      requestRef.current = requestAnimationFrame(processFrame)
    } catch (err) {
      console.error("Camera/MediaPipe initialization failed:", err)
      setHasPermission(false)
    }
  }

  const handleEyesLost = () => {
    setEyePositions(null)
    eyeDetectedRef.current = false
    if (capturePhase >= 0 && capturePhase < ritualPhrases.length) setCapturePhase(-1)
  }
  
  const renderFrame = (results: any) => {
    if (!canvasRef.current || !videoRef.current) return
    const ctx = canvasRef.current.getContext('2d')
    if (!ctx) return
    const cw = canvasRef.current.width
    const ch = canvasRef.current.height
    ctx.clearRect(0, 0, cw, ch)
    ctx.save()
    ctx.filter = 'blur(15px) brightness(0.3)'
    ctx.scale(-1, 1)
    ctx.drawImage(videoRef.current, -cw, 0, cw, ch)
    ctx.restore()
    if (results.multiFaceLandmarks && results.multiFaceLandmarks.length > 0) {
      const landmarks = results.multiFaceLandmarks[0]
      const leftEye = landmarks[LEFT_IRIS_CENTER]
      const rightEye = landmarks[RIGHT_IRIS_CENTER]
      if (leftEye && rightEye) {
        ctx.save()
        ctx.beginPath()
        const leftX = cw - (leftEye.x * cw)
        const leftY = leftEye.y * ch
        ctx.ellipse(leftX, leftY, 50, 35, 0, 0, Math.PI * 2)
        const rightX = cw - (rightEye.x * cw)
        const rightY = rightEye.y * ch
        ctx.ellipse(rightX, rightY, 50, 35, 0, 0, Math.PI * 2)
        ctx.clip()
        ctx.filter = 'none'
        ctx.scale(-1, 1)
        ctx.drawImage(videoRef.current, -cw, 0, cw, ch)
        ctx.restore()
        const currentPhase = capturePhaseRef.current
        const progress = currentPhase >= 0 ? (currentPhase + 1) / ritualPhrases.length : 0
        drawEyeProgressRing(ctx, leftEye, cw, ch, progress)
        drawEyeProgressRing(ctx, rightEye, cw, ch, progress)
      }
    }
  }
  
  const drawEyeProgressRing = (ctx: CanvasRenderingContext2D, eyeCenter: any, cw: number, ch: number, progress: number) => {
    if (!eyeCenter) return
    const x = cw - (eyeCenter.x * cw)
    const y = eyeCenter.y * ch
    const r = 17
    const lineWidth = 3
    const time = Date.now()
    const segments = 5
    const gapAngle = 0.08
    const segmentAngle = (Math.PI * 2 - gapAngle * segments) / segments
    const startOffset = -Math.PI / 2
    ctx.save()
    for (let i = 0; i < segments; i++) {
      const segmentStart = startOffset + i * (segmentAngle + gapAngle)
      const segmentEnd = segmentStart + segmentAngle
      const isLit = (i + 1) / segments <= progress
      const isCurrentSegment = i === Math.floor(progress * segments) && progress < 1
      ctx.beginPath()
      ctx.arc(x, y, r, segmentStart, segmentEnd)
      ctx.lineCap = 'round'
      ctx.lineWidth = lineWidth
      if (isLit) {
        const gradient = ctx.createLinearGradient(x + Math.cos(segmentStart) * r, y + Math.sin(segmentStart) * r, x + Math.cos(segmentEnd) * r, y + Math.sin(segmentEnd) * r)
        gradient.addColorStop(0, '#00ff88')
        gradient.addColorStop(1, '#00ffff')
        ctx.strokeStyle = gradient
        ctx.shadowBlur = 10
        ctx.shadowColor = '#00ffff'
        ctx.stroke()
      } else if (isCurrentSegment) {
        const segmentProgress = (progress * segments) % 1
        const partialEnd = segmentStart + segmentAngle * segmentProgress
        ctx.strokeStyle = 'rgba(0, 255, 255, 0.15)'
        ctx.stroke()
        ctx.beginPath()
        ctx.arc(x, y, r, segmentStart, partialEnd)
        ctx.strokeStyle = '#00ffff'
        ctx.stroke()
      } else {
        ctx.strokeStyle = 'rgba(0, 255, 255, 0.15)'
        ctx.stroke()
      }
    }
    const scanAngle = (time / 600) % (Math.PI * 2)
    ctx.strokeStyle = 'rgba(0, 255, 255, 0.8)'
    ctx.lineWidth = 2
    ctx.beginPath()
    ctx.moveTo(x + Math.cos(scanAngle) * (r - 20), y + Math.sin(scanAngle) * (r - 20))
    ctx.lineTo(x + Math.cos(scanAngle) * (r + 20), y + Math.sin(scanAngle) * (r + 20))
    ctx.stroke()
    ctx.restore()
  }

  const stopCamera = () => {
    isActiveRef.current = false
    if (requestRef.current) cancelAnimationFrame(requestRef.current)
    if (streamRef.current) streamRef.current.getTracks().forEach(track => track.stop())
    if (faceMeshRef.current) try { faceMeshRef.current.close() } catch(e) {}
    setIsStarted(false)
    setEyePositions(null)
    setLoadingStatus("IDLE")
    setCapturePhase(-1)
    setCapturedImage(null)
    setSoulCode(null)
    eyeDetectedRef.current = false
  }

  useEffect(() => { return () => stopCamera() }, [])

  if (hasPermission === false) return (
    <div className="flex flex-col items-center justify-center p-8 border border-black/20 bg-black/5 text-center">
      <Camera className="w-8 h-8 mb-3 opacity-20" />
      <p className="text-xs text-black font-normal font-sans">需要摄像头权限来捕捉你的灵魂印记</p>
      <button onClick={startCamera} className="mt-4 px-4 py-2 text-[10px] tracking-widest border border-black hover:bg-black hover:text-white transition-colors uppercase font-sans">启用摄像头</button>
    </div>
  )

  if (!isStarted) return (
    <button onClick={startCamera} className="w-full aspect-square flex flex-col items-center justify-center border border-black/20 bg-black/[0.02] hover:bg-black/[0.05] transition-all group">
      <div className="relative mb-4"><Scan className="w-10 h-10 opacity-20 group-hover:opacity-40 transition-opacity" /><div className="absolute inset-0 border border-black/20 scale-150 rounded-full animate-pulse" /></div>
      <span className="text-[10px] tracking-[0.2em] text-black font-normal uppercase font-sans">启动灵魂共振</span>
    </button>
  )

  const isReady = loadingStatus === "READY"
  if (capturedImage && capturePhase === 5) {
    const leftEyeX = 35, leftEyeY = 45, rightEyeX = 65, rightEyeY = 45
    return (
      <div className="relative aspect-square overflow-hidden border border-black bg-black">
        <img src={capturedImage} alt="Captured" className="absolute inset-0 w-full h-full object-cover blur-md opacity-50" />
        <div className="absolute inset-0 flex flex-col items-center justify-center pointer-events-none">
          <div className="text-center px-6 animate-pulse">
            <p className="text-[10px] uppercase tracking-[0.3em] text-[#00ffff]/80 mb-2 font-sans">灵魂印记已捕获</p>
            <p className="text-[#00ffff] font-mono text-sm tracking-wider mb-4">{soulCode}</p>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="relative aspect-square overflow-hidden border border-black bg-black">
      <video ref={videoRef} className="hidden" playsInline muted autoPlay />
      <canvas ref={canvasRef} width={480} height={480} className="w-full h-full object-cover" />
      {!isReady && <div className="absolute inset-0 bg-black flex flex-col items-center justify-center gap-4 z-20"><div className="relative w-16 h-16"><div className="absolute inset-0 border-2 border-white/10 rounded-full"></div><div className="absolute inset-0 border-2 border-white border-t-transparent rounded-full animate-spin"></div></div><p className="text-white/60 font-mono text-[10px] tracking-widest animate-pulse font-sans">{loadingStatus}</p></div>}
      {isReady && <div className="absolute inset-0 pointer-events-none flex flex-col"><div className="absolute top-4 left-4 flex items-center gap-2"><span className={`w-2 h-2 rounded-full ${eyePositions ? 'bg-green-500 animate-pulse' : 'bg-amber-500'}`}></span><span className="text-[9px] text-white/60 font-mono uppercase tracking-widest font-sans">{eyePositions ? '虹膜锁定' : '搜索中'}</span></div><div className="flex-1 flex items-end justify-center pb-16 px-6"><div className="text-center max-w-xs">{eyePositions && capturePhase >= 0 ? <p className="text-white/90 text-sm font-light italic leading-relaxed animate-fade-in font-serif">{ritualPhrases[capturePhase]}</p> : <p className="text-white/60 text-xs font-light leading-relaxed font-serif">{waitingMessage}</p>}</div></div></div>}
      {isReady && <button onClick={stopCamera} className="absolute top-4 right-4 p-2 bg-white/50 backdrop-blur rounded-full hover:bg-white transition-colors pointer-events-auto"><RefreshCw className="w-4 h-4" /></button>}
    </div>
  )
}

function OrbitalNode({ config, position, value, onSelect }: { config: typeof nodeConfigs[0], position: { top: string; left: string }, value: string, onSelect: () => void }) {
  const filled = Boolean(value)
  const display = value || config.label
  let iconName = config.icon
  let ZodiacIcon: any = null
  if (config.field === "soul" && value) iconName = "visibility"
  else if (config.field === "gender" && value) { if (value === "女") iconName = "female"; else if (value === "男") iconName = "male"; else iconName = "transgender"; }
  else if (config.field === "birthday" && value) { const zodiac = getZodiacFromBirthday(value); if (zodiac) ZodiacIcon = zodiac.icon; }
  return (
    <div className="absolute transition-all duration-500 ease-out" style={{ top: position.top, left: position.left, transform: "translate(-50%, -50%)" }}>
      <button type="button" onClick={onSelect} className="flex flex-col items-center text-[10px] tracking-[0.18em] transition-all duration-500 ease-out group">
        <div className={cn("flex items-center justify-center w-[56px] h-[56px] md:w-[64px] md:h-[64px] rounded-full border border-black dark:border-white bg-white dark:bg-black shadow-sm transition-all duration-300", filled ? "bg-black text-white dark:bg-white dark:text-black" : "hover:bg-black hover:text-white dark:hover:bg-white dark:hover:text-black")}>
          {ZodiacIcon ? <ZodiacIcon className="w-6 h-6" /> : <MaterialIcon name={iconName} className={cn(filled ? "" : "opacity-60")} />}
        </div>
        <span className={cn("absolute top-full mt-1.5 text-[10px] tracking-[0.04em] font-normal whitespace-nowrap max-w-[70px] truncate transition-colors font-sans uppercase", filled ? "text-black dark:text-white" : "text-black/50 dark:text-white/50")}>{config.field === "soul" && value ? "已捕获" : display}</span>
      </button>
    </div>
  )
}

export function ConnectUnified({ isVisible, onAccept, onReject, onComplete, onExit }: ConnectUnifiedProps) {
  const [phase, setPhase] = useState<Phase>("choice")
  const [animationState, setAnimationState] = useState<'hidden' | 'visible' | 'exiting'>('hidden')
  const [selectionStep, setSelectionStep] = useState<SelectionStep>("first")
  const [currentUserIndex, setCurrentUserIndex] = useState(0)
  const [firstSelectedUser, setFirstSelectedUser] = useState<UserOrb | null>(null)
  const [users, setUsers] = useState<UserOrb[]>(mockUsers)
  const [slideDirection, setSlideDirection] = useState<'left' | 'right' | null>(null)
  const [editingUser, setEditingUser] = useState<UserOrb | null>(null)
  const [newUserValues, setNewUserValues] = useState<Record<string, string>>({ name: "", gender: "", birthday: "", soul: "", blood: "" })
  const [editingField, setEditingField] = useState<string | null>(null)

  useEffect(() => {
    if (isVisible) {
      setPhase("choice"); setSelectionStep("first"); setCurrentUserIndex(0); setFirstSelectedUser(null);
      setNewUserValues({ name: "", gender: "", birthday: "", soul: "", blood: "" }); setEditingField(null);
      setAnimationState('hidden'); setTimeout(() => setAnimationState('visible'), 50);
    }
  }, [isVisible])

  const handlePrevUser = useCallback(() => { setSlideDirection('right'); setTimeout(() => { setCurrentUserIndex((prev) => (prev - 1 + users.length) % users.length); setSlideDirection(null); }, 150); }, [users.length])
  const handleNextUser = useCallback(() => { setSlideDirection('left'); setTimeout(() => { setCurrentUserIndex((prev) => (prev + 1) % users.length); setSlideDirection(null); }, 150); }, [users.length])
  const handleContinue = useCallback(() => {
    if (selectionStep === "first") { setFirstSelectedUser(users[currentUserIndex]); setSelectionStep("second"); setCurrentUserIndex((prev) => (prev + 1) % users.length); }
    else { const secondUser = users[currentUserIndex]; setPhase("resonating"); onAccept(); setTimeout(() => { setPhase("doorOpening"); setTimeout(() => { setPhase("result"); onComplete({ user1: firstSelectedUser, user2: secondUser }); }, 2000); }, 6000); }
  }, [selectionStep, users, currentUserIndex, firstSelectedUser, onAccept, onComplete])
  const handleCancel = useCallback(() => { if (selectionStep === "second") { setSelectionStep("first"); if (firstSelectedUser) { const idx = users.findIndex(u => u.id === firstSelectedUser.id); if (idx !== -1) setCurrentUserIndex(idx); } } else { setAnimationState('exiting'); setTimeout(() => onReject(), 500); } }, [selectionStep, firstSelectedUser, users, onReject])
  const handleNewUser = useCallback(() => { setNewUserValues({ name: "", gender: "", birthday: "", soul: "", blood: "" }); setEditingField(null); setEditingUser(null); setPhase("creating"); }, [])
  const handleEditUser = useCallback(() => { const user = users[currentUserIndex]; setEditingUser(user); setNewUserValues({ name: user.name, gender: user.gender, birthday: user.birthday, soul: user.soul, blood: user.blood }); setEditingField(null); setPhase("editing"); }, [users, currentUserIndex])
  const handleConfirmNewUser = useCallback(() => { if (!newUserValues.name.trim()) return; const newUser: UserOrb = { id: `new-${Date.now()}`, name: newUserValues.name.trim(), gender: newUserValues.gender, birthday: newUserValues.birthday, soul: newUserValues.soul, blood: newUserValues.blood }; setUsers(prev => [...prev, newUser]); setCurrentUserIndex(users.length); setPhase("choice"); }, [newUserValues, users.length])
  const handleCancelNewUser = useCallback(() => { setPhase("choice"); setEditingField(null); setEditingUser(null); }, [])
  const handleSaveEditUser = useCallback(() => { if (!editingUser || !newUserValues.name.trim()) return; setUsers(prev => prev.map(u => u.id === editingUser.id ? { ...u, ...newUserValues } : u)); setPhase("choice"); setEditingUser(null); }, [editingUser, newUserValues])
  const handleDeleteUser = useCallback(() => { if (!editingUser) return; setUsers(prev => { const filtered = prev.filter(u => u.id !== editingUser.id); if (currentUserIndex >= filtered.length) setCurrentUserIndex(Math.max(0, filtered.length - 1)); return filtered; }); setPhase("choice"); setEditingUser(null); }, [editingUser, currentUserIndex])
  const handleSelectField = useCallback((field: string) => setEditingField(field), [])
  const handleExit = useCallback(() => { setPhase("exiting"); setTimeout(() => onExit(), 500); }, [onExit])

  if (!isVisible) return null
  const isChoicePhase = phase === "choice", isCreatingPhase = phase === "creating", isEditingPhase = phase === "editing", isResonating = phase === "resonating", isDoorOpening = phase === "doorOpening", isResultPhase = phase === "result"
  const secondSelectedUser = selectionStep === "second" ? users[currentUserIndex] : null

  return (
    <div className={`fixed z-[100] inset-0 transition-all duration-1000 ${phase === "exiting" ? "opacity-0" : "opacity-100"} pointer-events-none`}>
      <div className={cn("absolute left-0 right-0 transition-all duration-[1200ms] cubic-bezier(0.4, 0, 0.2, 1)", isChoicePhase ? "bg-white/80 dark:bg-black/80 backdrop-blur-md" : "bg-white dark:bg-black")} style={{ bottom: isChoicePhase ? "64px" : "0", top: isChoicePhase ? "calc(100% - 384px)" : "0", pointerEvents: isChoicePhase ? "none" : "auto" }} />
      {(phase === "interaction" || phase === "result") && <button onClick={handleExit} className="absolute top-6 right-6 z-50 p-2 border border-black/30 dark:border-white/30 hover:bg-black/10 dark:hover:bg-white/10 transition-colors pointer-events-auto"><X className="w-5 h-5 text-black dark:text-white" /></button>}
      {isChoicePhase && (
        <div className="absolute inset-0 pointer-events-none flex flex-col items-center justify-center" style={{ transform: "translateY(calc(50vh - 224px - 40px))", opacity: animationState === 'visible' ? 1 : 0, filter: animationState === 'visible' ? 'blur(0px)' : 'blur(10px)', transition: "transform 1.2s cubic-bezier(0.16, 1, 0.3, 1), opacity 0.8s ease-out, filter 0.8s ease-out" }}>
          <div className="relative flex items-center gap-6 pointer-events-auto">
            <button onClick={handlePrevUser} className="p-2 hover:bg-black/10 dark:hover:bg-white/10 transition-colors rounded-full"><ChevronLeft className="w-6 h-6 text-black dark:text-white" /></button>
            <div className="relative">
              {users[currentUserIndex] && <div className="absolute -bottom-3 left-1/2 -translate-x-1/2 z-10 px-3 py-1 bg-black dark:bg-white text-white dark:text-black text-xs font-medium shadow-md whitespace-nowrap font-sans uppercase" style={{ opacity: slideDirection ? 0 : 1, transition: "opacity 0.15s ease-out" }}>{users[currentUserIndex].name}</div>}
              <button onClick={handleEditUser} className="relative rounded-full w-[100px] h-[100px] bg-white dark:bg-black border border-black dark:border-white shadow-sm animate-[spinSlow_18s_linear_infinite] cursor-pointer" style={{ opacity: slideDirection ? 0.5 : 1, transform: slideDirection === 'left' ? 'translateX(-20px)' : slideDirection === 'right' ? 'translateX(20px)' : 'translateX(0)', transition: "opacity 0.15s ease-out, transform 0.15s ease-out" }}>
                <svg viewBox="0 0 100 100" className="absolute inset-0 text-black dark:text-white">
                  <circle cx="50" cy="50" r="46" fill="none" stroke="currentColor" strokeWidth={0.5} strokeDasharray="2 4" />
                  <polygon points={starPoints} fill="none" stroke="currentColor" strokeWidth={0.5} strokeLinecap="round" strokeLinejoin="round" />
                  {ringPositions.map((pos, i) => { const fieldKeys = ["name", "gender", "birthday", "soul", "blood"] as const; const filled = users[currentUserIndex] && Boolean(users[currentUserIndex][fieldKeys[i]]); return filled ? <circle key={i} cx={pos.x} cy={pos.y} r="3" fill="currentColor" /> : <circle key={i} cx={pos.x} cy={pos.y} r="2.5" fill="none" stroke="currentColor" strokeWidth={0.5} /> })}
                </svg>
              </button>
            </div>
            <button onClick={handleNextUser} className="p-2 hover:bg-black/10 dark:hover:bg-white/10 transition-colors rounded-full"><ChevronRight className="w-6 h-6 text-black dark:text-white" /></button>
          </div>
          <div className="absolute w-full flex flex-col items-center pointer-events-none" style={{ top: "calc(50% + 80px)" }}>
             <h3 className="text-lg font-light tracking-wide mb-6 font-serif text-black dark:text-white"><InkRevealText key={selectionStep} text={selectionStep === "first" ? "请选择第一个用户" : "请选择第二个用户"} /></h3>
             <div className="flex gap-3 w-full justify-center px-6 max-w-md pointer-events-auto">
               <button onClick={handleContinue} className="flex-1 max-w-[100px] border border-black dark:border-white py-3 text-sm font-light bg-black dark:bg-white text-white dark:text-black hover:bg-white dark:hover:bg-black hover:text-black dark:hover:text-white transition-colors shadow-sm font-sans uppercase">继续</button>
               <button onClick={handleCancel} className="flex-1 max-w-[100px] border border-black dark:border-white py-3 text-sm font-light bg-white dark:bg-black text-black dark:text-white hover:bg-black hover:text-white transition-colors shadow-sm font-sans uppercase">{selectionStep === "second" ? "返回" : "取消"}</button>
               <button onClick={handleNewUser} className="flex-1 max-w-[120px] border border-black dark:border-white py-3 text-sm font-light bg-white dark:bg-black text-black dark:text-white hover:bg-black hover:text-white transition-colors shadow-sm flex items-center justify-center gap-1 font-sans uppercase"><Plus className="w-4 h-4" />新建用户</button>
             </div>
          </div>
        </div>
      )}
      {(isCreatingPhase || isEditingPhase) && (
        <div className="absolute inset-0 bg-white dark:bg-black pointer-events-auto flex flex-col items-center">
          <div className="pt-16 pb-6"><h2 className="text-xl font-light tracking-wide font-serif text-black dark:text-white"><InkRevealText text={isCreatingPhase ? "添加用户信息" : "修改用户信息"} /></h2></div>
          <div className="flex-1 flex items-center justify-center">
            <div className="relative rounded-full w-[280px] h-[280px] bg-white dark:bg-black border border-black dark:border-white shadow-sm">
              <svg viewBox="0 0 100 100" className="absolute inset-0 text-black/30 dark:text-white/30"><circle cx="50" cy="50" r="46" fill="none" stroke="currentColor" strokeWidth={0.5} strokeDasharray="2 4" /><polygon points={starPoints} fill="none" stroke="currentColor" strokeWidth={0.3} strokeDasharray="1 2" strokeLinecap="round" strokeLinejoin="round" /></svg>
              {nodeConfigs.map((config, index) => <OrbitalNode key={config.field} config={config} position={ringPositions[index]} value={newUserValues[config.field]} onSelect={() => handleSelectField(config.field)} />)}
            </div>
          </div>
          <div className="pb-16 pt-8 flex flex-col gap-4 items-center w-full px-6">
            <div className="flex gap-3 w-full max-w-xs">
              <button onClick={isCreatingPhase ? handleConfirmNewUser : handleSaveEditUser} disabled={!newUserValues.name.trim()} className="flex-1 py-3.5 border border-black dark:border-white text-sm font-light bg-black dark:bg-white text-white dark:text-black hover:bg-white dark:hover:bg-black hover:text-black dark:hover:text-white transition-colors disabled:opacity-40 font-sans uppercase">确认{isCreatingPhase ? '添加' : '修改'}</button>
              <button onClick={handleCancelNewUser} className="flex-1 py-3.5 border border-black dark:border-white text-sm font-light bg-white dark:bg-black text-black dark:text-white hover:bg-black hover:text-white transition-colors font-sans uppercase">取消</button>
            </div>
            {isEditingPhase && <button onClick={handleDeleteUser} className="text-xs tracking-wider text-black/40 dark:text-white/40 hover:text-red-500 transition-colors font-sans uppercase">删除此用户</button>}
          </div>
          {editingField && (
            <>
              <div className="fixed inset-0 z-30 backdrop-blur-md bg-white/50 dark:bg-black/50" onClick={() => setEditingField(null)} />
              <div className="fixed left-1/2 -translate-x-1/2 bottom-24 w-[92%] max-w-md z-40" onClick={(e) => e.stopPropagation()}>
                <div className="border border-black/50 dark:border-white/50 bg-white dark:bg-black backdrop-blur p-5 shadow-xl">
                  <div className="flex items-start justify-between gap-3 mb-4"><div><p className="text-sm tracking-wide text-black dark:text-white font-normal font-sans uppercase">{nodeConfigs.find(c => c.field === editingField)?.label}</p></div><button type="button" className="text-xs text-black dark:text-white font-normal hover:opacity-70 font-sans uppercase" onClick={() => setEditingField(null)}>关闭</button></div>
                  {editingField === "gender" ? (
                    <div className="flex flex-wrap gap-2">{["女", "男", "其他"].map((option) => <button key={option} type="button" onClick={() => { setNewUserValues(prev => ({ ...prev, gender: option })); setEditingField(null); }} className={cn("px-4 py-2 text-sm border font-normal transition-colors font-sans", newUserValues.gender === option ? "bg-black text-white dark:bg-white dark:text-black" : "bg-white dark:bg-black text-black dark:text-white hover:bg-black hover:text-white dark:hover:bg-white dark:hover:text-black")}>{option}</button>)}</div>
                  ) : editingField === "blood" ? (
                    <div className="flex flex-wrap gap-2">{["A", "B", "AB", "O"].map((option) => <button key={option} type="button" onClick={() => { setNewUserValues(prev => ({ ...prev, blood: option })); setEditingField(null); }} className={cn("px-4 py-2 text-sm border font-normal transition-colors font-sans", newUserValues.blood === option ? "bg-black text-white dark:bg-white dark:text-black" : "bg-white dark:bg-black text-black dark:text-white hover:bg-black hover:text-white dark:hover:bg-white dark:hover:text-black")}>{option}</button>)}</div>
                  ) : editingField === "birthday" ? (
                    <div className="flex items-center border border-black dark:border-white"><span className="px-3 text-xs text-black dark:text-white font-normal font-sans uppercase">日期</span><input type="date" value={newUserValues.birthday} onChange={(e) => setNewUserValues(prev => ({ ...prev, birthday: e.target.value }))} className="flex-1 px-2 py-2.5 text-sm bg-transparent text-black dark:text-white font-normal focus:outline-none" /></div>
                  ) : editingField === "soul" ? (
                    <div className="flex flex-col gap-4">{newUserValues.soul ? <div className="flex flex-col"><div className="flex items-center justify-center gap-2 py-4"><CheckCircle2 className="w-5 h-5 text-black dark:text-white" /><span className="text-sm font-mono tracking-wider text-black dark:text-white">{newUserValues.soul}</span></div></div> : <EyeScanner onCapture={(code) => { setNewUserValues(prev => ({ ...prev, soul: code })); }} />}</div>
                  ) : (
                    <input type="text" value={newUserValues[editingField] || ""} onChange={(e) => setNewUserValues(prev => ({ ...prev, [editingField]: e.target.value }))} placeholder={`输入${nodeConfigs.find(c => c.field === editingField)?.label}`} className="w-full border border-black dark:border-white px-4 py-3 text-sm bg-white dark:bg-black text-black dark:text-white focus:outline-none font-sans" autoFocus />
                  )}
                  {editingField !== "gender" && editingField !== "blood" && (
                    <div className="flex justify-between items-center mt-4">
                      {editingField === "soul" && newUserValues.soul ? <button onClick={() => { setNewUserValues(prev => ({ ...prev, soul: "" })); }} className="text-[10px] tracking-widest text-black/60 dark:text-white/60 hover:text-black dark:hover:text-white font-sans uppercase">重新扫描</button> : <div></div>}
                      <div className="flex gap-2"><button type="button" className="px-4 py-2 text-sm border border-black/40 dark:border-white/40 text-black dark:text-white hover:bg-black hover:text-white transition-colors font-sans uppercase" onClick={() => setEditingField(null)}>取消</button><button type="button" className="px-4 py-2 text-sm border border-black dark:border-white bg-black dark:bg-white text-white dark:text-black hover:bg-white hover:text-black transition-colors font-sans uppercase" onClick={() => setEditingField(null)}>保存</button></div>
                    </div>
                  )}
                </div>
              </div>
            </>
          )}
        </div>
      )}
      {isResonating && (
        <div className="fixed inset-0 z-[100] pointer-events-none">
          <div className="absolute inset-0 z-10 transition-opacity duration-1000 opacity-95 bg-white dark:bg-black" />
          <div className="absolute inset-0 z-20 flex items-center justify-center">
            <section className="relative flex flex-col items-center gap-10">
              <div className="relative z-30" style={{ animation: "portalOrbit 6s ease-in-out forwards", transformOrigin: "center calc(100% + 40px)" }}>
                <div className="relative"><div className="absolute -inset-1 rounded-full z-[-1] laser-gradient opacity-50 blur-md" style={{ animation: 'spinSlow 1.5s linear infinite' }} /><div className="relative rounded-full w-[100px] h-[100px] bg-white dark:bg-black border border-black dark:border-white shadow-sm animate-[spinSlow_18s_linear_infinite]"><svg viewBox="0 0 100 100" className="absolute inset-0 text-black dark:text-white"><circle cx="50" cy="50" r="46" fill="none" stroke="currentColor" strokeWidth={0.5} strokeDasharray="2 4" /><polygon points={starPoints} fill="none" stroke="currentColor" strokeWidth={0.5} strokeLinecap="round" strokeLinejoin="round" />{ringPositions.map((pos, i) => { const filled = firstSelectedUser && Boolean(firstSelectedUser[i === 0 ? "name" : i === 1 ? "gender" : i === 2 ? "birthday" : i === 3 ? "soul" : "blood"]); return filled ? <circle key={i} cx={pos.x} cy={pos.y} r="3" fill="currentColor" /> : <circle key={i} cx={pos.x} cy={pos.y} r="2.5" fill="none" stroke="currentColor" strokeWidth={0.5} /> })}</svg></div></div>
              </div>
              <div className="relative flex items-center justify-center w-full py-1"><div className="relative flex items-center justify-center w-16 h-16 rounded-full border border-black dark:border-white bg-white dark:bg-black z-50"><div className="absolute inset-[6px] rounded-full bg-black dark:bg-white flex items-center justify-center overflow-hidden"><OracleInteractiveEye direction="center" className="w-9 h-9" inverted={true} /></div></div></div>
              <div className="relative z-30" style={{ animation: "portalOrbit 6s ease-in-out forwards", transformOrigin: "center calc(0% - 40px)" }}>
                <div className="relative"><div className="absolute -inset-1 rounded-full z-[-1] laser-gradient opacity-50 blur-md" style={{ animation: 'spinSlow 2.1s linear infinite reverse' }} /><div className="relative rounded-full w-[100px] h-[100px] bg-white dark:bg-black border border-black dark:border-white shadow-sm animate-[spinSlow_18s_linear_infinite]"><svg viewBox="0 0 100 100" className="absolute inset-0 text-black dark:text-white"><circle cx="50" cy="50" r="46" fill="none" stroke="currentColor" strokeWidth={0.5} strokeDasharray="2 4" /><polygon points={starPoints} fill="none" stroke="currentColor" strokeWidth={0.5} strokeLinecap="round" strokeLinejoin="round" />{ringPositions.map((pos, i) => { const filled = secondSelectedUser && Boolean(secondSelectedUser[i === 0 ? "name" : i === 1 ? "gender" : i === 2 ? "birthday" : i === 3 ? "soul" : "blood"]); return filled ? <circle key={i} cx={pos.x} cy={pos.y} r="3" fill="currentColor" /> : <circle key={i} cx={pos.x} cy={pos.y} r="2.5" fill="none" stroke="currentColor" strokeWidth={0.5} /> })}</svg></div></div>
              </div>
            </section>
          </div>
          <div className="absolute inset-0 z-30 flex items-end justify-center pb-10"><p className="text-center text-black dark:text-white text-xs font-normal animate-pulse font-serif uppercase tracking-widest">正在校准灵魂共振频率...</p></div>
        </div>
      )}
      <div className={cn("fixed inset-0 z-[100] bg-black flex flex-col items-center transition-all duration-[1500ms] ease-in-out", (isDoorOpening || isResultPhase) ? "opacity-100 visible pointer-events-auto" : "opacity-0 invisible pointer-events-none")} style={{ clipPath: (isDoorOpening || isResultPhase) ? 'inset(0% 0% 0% 0%)' : 'inset(50% 0% 50% 0%)' }}>
        <SolarSystem active={isDoorOpening || isResultPhase} />
        {isResultPhase && (
          <>
            <button onClick={handleExit} className="absolute top-6 right-6 z-50 p-2 border border-white/30 hover:bg-white/10 transition-colors"><X className="w-5 h-5 text-white" /></button>
            <div className="absolute top-8 left-1/2 -translate-x-1/2 z-40"><div className="relative flex items-center justify-center w-16 h-16 rounded-full border border-white bg-black"><div className="absolute inset-[6px] rounded-full bg-white flex items-center justify-center overflow-hidden"><OracleInteractiveEye direction="center" className="w-9 h-9" inverted={false} /></div></div></div>
            <div className="absolute inset-0 z-10 flex flex-col items-center pt-[114px] animate-fade-in-delayed">
              <div className="w-full max-w-2xl px-8 text-center shrink-0 pb-6"><h2 className="text-xl md:text-2xl font-light text-white tracking-[0.3em] uppercase whitespace-nowrap font-serif"><InkRevealText text="Celestial Synchrony" /></h2></div>
              <div className="flex-1 w-full max-w-2xl px-8 overflow-y-auto hide-scrollbar">
                <div className="text-sm md:text-base text-white font-light leading-[1.8] text-justify space-y-8 pb-48 font-serif">
                  <p>在浩瀚的宇宙构架中，两个独特的光谱特征在精确的时间节点相交。这种相交不仅仅是物理领域的偶然相遇，而是穿越多个意识维度、到达这一统一清明时刻的振动状态的深刻对齐。{firstSelectedUser?.name} 与 {secondSelectedUser?.name} 之间观察到的和谐汇聚，指向一种先于个体存在的古老共鸣。</p>
                  <p>来自眼纹印记分析的数据流揭示了共享原型模式的复杂织锦。虹膜结构中的每一次微振动都充当着宇宙历史的生物记录，当这些记录同步时，它们会创造出可在量子场中测量的共振频率。这种对齐表明灵魂纠缠的高概率——意识进化的反馈循环。</p>
                  <p>这种对齐超越了个性特质或共同兴趣。它延伸到光环场的核心，在那里南北节点的基本能量极性达到罕见的平衡。这种平衡为共同目的提供了稳定的基础，使两位参与者的组合电荷能够以最小损失穿越世俗世界的不和谐频率。</p>
                  <p>在切实的层面上，这种共振表现为沟通中的轻松节奏、同步决策，以及相互支持感。当挑战出现时，组合场充当稳定器：它减缓反应性螺旋并增强清晰度。这不是摩擦的缺失，而是一种将摩擦转化为前进动力的底层对齐的存在。</p>
                  <p>最终的综合表明这种联系并非偶然。它作为系统更大和谐中的灯塔发挥作用。如果你们用诚实、时间和关注来滋养它，这种共振将保持为一种活的乐器：响应性的、进化的，并能够在不削弱任何一方自我的情况下引导你们穿越不确定性。让它深思熟虑。让它善良。让它真实。</p>
                </div>
              </div>
            </div>
            <div className="absolute bottom-8 left-0 right-0 z-30 flex justify-center gap-4 px-8"><button onClick={handleExit} className="px-8 py-3 border border-white/30 text-sm font-light text-white hover:bg-white/10 transition-colors font-sans uppercase">返回</button><button onClick={() => alert("分享功能开发中...")} className="px-8 py-3 border border-white text-sm font-light bg-white text-black hover:bg-white/90 transition-colors font-sans uppercase">分享到 Story</button></div>
            <div className="absolute inset-x-0 bottom-0 h-32 bg-gradient-to-t from-black via-black/80 to-transparent pointer-events-none z-20" />
          </>
        )}
      </div>
      <style>{`
        @keyframes spinSlow { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
        @keyframes portalOrbit { 0% { transform: rotate(0deg); } 100% { transform: rotate(2160deg); } }
        .animate-fade-in-delayed { animation: fadeIn 1s ease-out 0.5s forwards; opacity: 0; }
        @keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }
      `}</style>
    </div>
  )
}

