import { useEffect, useRef } from "react"
import * as THREE from "three"
import { RectAreaLightUniformsLib } from "three/examples/jsm/lights/RectAreaLightUniformsLib.js"

// 初始化 RectAreaLight 所需的 uniforms
RectAreaLightUniformsLib.init()

interface Stone3DProps {
  size?: number
  isHolding?: boolean
  progress?: number
  revealed?: boolean
  isShaking?: boolean
  onClick?: () => void
}

export function Stone3D({ 
  size = 256, 
  isHolding = false, 
  progress = 0, 
  revealed = false,
  isShaking = false,
  onClick 
}: Stone3DProps) {
  const containerRef = useRef<HTMLDivElement>(null)
  const sceneRef = useRef<{
    scene: THREE.Scene
    camera: THREE.PerspectiveCamera
    renderer: THREE.WebGLRenderer
    stone: THREE.Mesh
    lights: THREE.RectAreaLight[]
    animationId: number
  } | null>(null)
  
  // 使用 ref 跟踪状态，以便在动画循环中访问
  const isHoldingRef = useRef(isHolding)
  const progressRef = useRef(progress)
  const revealedRef = useRef(revealed)
  const isShakingRef = useRef(isShaking)
  
  useEffect(() => {
    isHoldingRef.current = isHolding
    progressRef.current = progress
  }, [isHolding, progress])
  
  useEffect(() => {
    revealedRef.current = revealed
  }, [revealed])
  
  useEffect(() => {
    isShakingRef.current = isShaking
  }, [isShaking])

  useEffect(() => {
    if (!containerRef.current) return

    const container = containerRef.current

    // 创建场景
    const scene = new THREE.Scene()

    // 创建相机 - 拉远确保石头放大时不会超出画布
    const camera = new THREE.PerspectiveCamera(50, 1, 0.1, 1000)
    camera.position.z = 5

    // 创建渲染器
    const renderer = new THREE.WebGLRenderer({ 
      antialias: true, 
      alpha: true 
    })
    renderer.setSize(size, size)
    renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2))
    renderer.toneMapping = THREE.ACESFilmicToneMapping
    renderer.toneMappingExposure = 1.2
    container.appendChild(renderer.domElement)

    // 创建二十面体（Icosahedron）作为多面体石头
    const geometry = new THREE.IcosahedronGeometry(1.5, 0)
    
    // 创建镜面反光材质 - 深灰色金属质感，类似黑曜石
    const material = new THREE.MeshPhysicalMaterial({
      color: 0x111111,
      metalness: 0.95,
      roughness: 0.02,
      reflectivity: 1,
      clearcoat: 1,
      clearcoatRoughness: 0.01,
      envMapIntensity: 5,
      ior: 2.5, // 高折射率增加反光
    })

    const stone = new THREE.Mesh(geometry, material)
    scene.add(stone)

    // 环境光 - 基础照明
    const ambientLight = new THREE.AmbientLight(0xffffff, 0.3)
    scene.add(ambientLight)

    // 正面大面积矩形面光源 - 模拟摄影棚柔光箱
    const rectLight = new THREE.RectAreaLight(0xffffff, 8, 10, 10)
    rectLight.position.set(0, 0, 5)
    rectLight.lookAt(0, 0, 0)
    scene.add(rectLight)
    
    // 辅助面光源 - 上方
    const rectLight2 = new THREE.RectAreaLight(0xffffff, 3, 8, 8)
    rectLight2.position.set(0, 5, 3)
    rectLight2.lookAt(0, 0, 0)
    scene.add(rectLight2)
    
    // 辅助面光源 - 侧面补光
    const rectLight3 = new THREE.RectAreaLight(0xffffff, 2, 6, 6)
    rectLight3.position.set(4, 0, 3)
    rectLight3.lookAt(0, 0, 0)
    scene.add(rectLight3)
    
    const rectLight4 = new THREE.RectAreaLight(0xffffff, 2, 6, 6)
    rectLight4.position.set(-4, 0, 3)
    rectLight4.lookAt(0, 0, 0)
    scene.add(rectLight4)

    // 创建简单的环境贴图
    const cubeRenderTarget = new THREE.WebGLCubeRenderTarget(256)
    const cubeCamera = new THREE.CubeCamera(0.1, 10, cubeRenderTarget)
    
    // 创建白色渐变背景球体用于环境反射（仅用于生成环境贴图，不渲染到屏幕）
    const envGeometry = new THREE.SphereGeometry(50, 32, 32)
    const envMaterial = new THREE.MeshBasicMaterial({
      side: THREE.BackSide,
      color: 0xffffff,
    })
    const envSphere = new THREE.Mesh(envGeometry, envMaterial)
    scene.add(envSphere)

    // 更新环境贴图
    stone.visible = false
    cubeCamera.position.copy(stone.position)
    cubeCamera.update(renderer, scene)
    stone.visible = true
    material.envMap = cubeRenderTarget.texture
    
    // 生成完环境贴图后移除背景球体，保持透明背景
    scene.remove(envSphere)

    // 存储引用（包括所有面光源）
    const lights = [rectLight, rectLight2, rectLight3, rectLight4]
    sceneRef.current = {
      scene,
      camera,
      renderer,
      stone,
      lights,
      animationId: 0,
    }

    // 动画循环
    let time = 0
    let colorTime = 0
    let currentOpacity = 1
    const animate = () => {
      if (!sceneRef.current) return
      
      const currentIsHolding = isHoldingRef.current
      const currentProgress = progressRef.current
      // const currentRevealed = revealedRef.current
      
      // 根据长按状态调整旋转速度
      const speedMultiplier = currentIsHolding ? 1 + currentProgress * 8 : 1 // 最高9倍速，更明显的加速
      time += 0.01 * speedMultiplier

      // 自动旋转
      stone.rotation.x = time * 0.3
      stone.rotation.y = time * 0.5
      
      // 震动效果
      const currentIsShaking = isShakingRef.current
      if (currentIsShaking) {
        stone.position.x = (Math.random() - 0.5) * 0.15
        stone.position.y = (Math.random() - 0.5) * 0.15
      } else {
        // 平滑恢复到中心
        stone.position.x *= 0.9
        stone.position.y *= 0.9
      }

      // 长按时石头变大
      const targetScale = currentIsHolding ? 1 + currentProgress * 0.5 : 1 // 最大1.5倍，更明显的放大
      stone.scale.lerp(new THREE.Vector3(targetScale, targetScale, targetScale), 0.1)

      // 长按时光源颜色变化
      if (currentIsHolding) {
        colorTime += 0.05
        const hue = (colorTime % 1) // 0-1 循环
        const color = new THREE.Color().setHSL(hue, 0.8, 0.6)
        lights.forEach(light => {
          light.color = color
        })
      } else {
        // 恢复白色光源
        colorTime = 0
        lights.forEach(light => {
          light.color.lerp(new THREE.Color(0xffffff), 0.1)
        })
      }
      
      // 保持石头不透明，不再受 revealed 影响
      const targetOpacity = 1
      currentOpacity += (targetOpacity - currentOpacity) * 0.08
      material.opacity = currentOpacity
      material.transparent = currentOpacity < 0.99

      renderer.render(scene, camera)
      sceneRef.current!.animationId = requestAnimationFrame(animate)
    }
    animate()

    // 清理
    return () => {
      if (sceneRef.current) {
        cancelAnimationFrame(sceneRef.current.animationId)
        sceneRef.current.renderer.dispose()
        container.removeChild(sceneRef.current.renderer.domElement)
        sceneRef.current = null
      }
    }
  }, [size])

  return (
    <div 
      ref={containerRef}
      onClick={onClick}
      className="cursor-pointer"
      style={{ 
        width: size, 
        height: size,
      }}
    />
  )
}

