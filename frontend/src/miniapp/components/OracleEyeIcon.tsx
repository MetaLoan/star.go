export function OracleEyeIcon({ className = "" }: { className?: string }) {
  // 尖锐眼眶的路径
  const eyePath = "M 10 100 Q 100 45 190 100 Q 100 155 10 100 Z"

  return (
    <div className={`flex items-center justify-center bg-black rounded-full ${className}`}>
      <svg
        width="100%"
        height="100%"
        viewBox="0 0 200 200"
        className="p-2"
      >
        <defs>
          {/* 尖锐眼眶的裁剪区域 */}
          <clipPath id="sharpEyeClip">
            <path d={eyePath} />
          </clipPath>
        </defs>

        {/* 装饰性外部线条 */}
        <path 
          d="M 0 100 L 30 100 M 170 100 L 200 100" 
          stroke="white" 
          strokeWidth="0.5" 
          opacity="0.3"
        />
        
        {/* 主眼眶轮廓 */}
        <path 
          d={eyePath} 
          fill="none" 
          stroke="white" 
          strokeWidth="1.2" 
        />
        
        {/* 内部细轮廓线 */}
        <path 
          d="M 18 100 Q 100 52 182 100 Q 100 148 18 100 Z" 
          fill="none" 
          stroke="white" 
          strokeWidth="0.5" 
          opacity="0.4"
        />

        {/* 圆形虹膜 */}
        <g clipPath="url(#sharpEyeClip)">
          <g>
            {/* 虹膜轮廓 */}
            <circle cx="100" cy="100" r="42" fill="none" stroke="white" strokeWidth="0.8" />
            <circle cx="100" cy="100" r="38" fill="none" stroke="white" strokeWidth="0.3" strokeDasharray="2 2" />
            
            {/* 放射状精密线条 */}
            {[...Array(48)].map((_, i) => (
              <line
                key={i}
                x1="100" y1="100"
                x2={Math.round((100 + Math.cos(i * 7.5 * Math.PI / 180) * 42) * 100) / 100}
                y2={Math.round((100 + Math.sin(i * 7.5 * Math.PI / 180) * 42) * 100) / 100}
                stroke="white"
                strokeWidth="0.3"
                opacity="0.2"
              />
            ))}

            {/* 核心瞳孔 - 竖直且尖锐 */}
            <ellipse 
              cx="100" cy="100" 
              rx="8" ry="28" 
              fill="white" 
            />
            
            {/* 瞳孔内部的黑色细缝 */}
            <ellipse 
              cx="100" cy="100" 
              rx="1.5" ry="20" 
              fill="black" 
            />
          </g>
        </g>

        {/* 眼角排线装饰 */}
        <g stroke="white" strokeWidth="0.5" opacity="0.2">
          <line x1="20" y1="95" x2="35" y2="85" />
          <line x1="20" y1="105" x2="35" y2="115" />
          <line x1="180" y1="95" x2="165" y2="85" />
          <line x1="180" y1="105" x2="165" y2="115" />
        </g>
      </svg>
    </div>
  )
}

