/**
 * LunarInfo7934IJ - 月亮信息组件
 * 显示月亮空亡、月相、行星时等信息
 */

import { motion } from 'framer-motion';

interface VoidOfCourseInfo {
  isVoid: boolean;
  duration: number;
  nextSign: string;
  lastAspect: string;
  influence: number;
}

interface PlanetaryHourInfo {
  planetaryHour: number;
  ruler: string;
  planetName: string;
  planetSymbol: string;
  dayRuler: string;
  influence: number;
  bestFor: string[];
}

interface LunarInfoProps {
  voidOfCourse?: VoidOfCourseInfo;
  planetaryHour?: PlanetaryHourInfo;
  moonPhase?: {
    name: string;
    phase: string;
    illumination: number;
  };
  moonSign?: string;
  className?: string;
}

export function LunarInfo7934IJ({
  voidOfCourse,
  planetaryHour,
  moonPhase,
  moonSign,
  className = '',
}: LunarInfoProps) {
  return (
    <div className={`glass-card p-4 ${className}`}>
      <h3 className="text-lg font-bold mb-4 flex items-center gap-2">
        🌙 月亮与时机
      </h3>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {/* 月亮空亡 */}
        {voidOfCourse && (
          <VoidOfCourseCard voidInfo={voidOfCourse} />
        )}

        {/* 行星时 */}
        {planetaryHour && (
          <PlanetaryHourCard hourInfo={planetaryHour} />
        )}
      </div>

      {/* 月相信息 */}
      {moonPhase && (
        <div className="mt-4 p-3 rounded-lg bg-white/5">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <MoonPhaseIcon phase={moonPhase.phase} />
              <div>
                <div className="font-medium">{moonPhase.name}</div>
                <div className="text-sm text-celestial-silver/60">
                  {moonSign && `月亮在${moonSign}`}
                </div>
              </div>
            </div>
            <div className="text-right">
              <div className="text-sm text-celestial-silver/60">光照度</div>
              <div className="text-lg font-bold">
                {(moonPhase.illumination * 100).toFixed(0)}%
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

// 月亮空亡卡片
function VoidOfCourseCard({ voidInfo }: { voidInfo: VoidOfCourseInfo }) {
  return (
    <motion.div
      className={`p-4 rounded-lg ${
        voidInfo.isVoid
          ? 'bg-cosmic-aurora/10 border border-cosmic-aurora/30'
          : 'bg-white/5'
      }`}
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
    >
      <div className="flex items-center justify-between mb-2">
        <div className="font-medium flex items-center gap-2">
          {voidInfo.isVoid ? '⚠️ 月亮空亡中' : '✓ 月亮活跃'}
        </div>
        {voidInfo.isVoid && (
          <div className="text-sm text-cosmic-aurora font-bold">
            {voidInfo.influence.toFixed(1)}
          </div>
        )}
      </div>

      {voidInfo.isVoid ? (
        <>
          <div className="text-sm text-celestial-silver/80 mb-2">
            持续 {voidInfo.duration.toFixed(1)} 小时后进入{voidInfo.nextSign}
          </div>
          <div className="text-xs text-celestial-silver/60">
            最后相位: {voidInfo.lastAspect}
          </div>
          <div className="mt-3 p-2 rounded bg-cosmic-aurora/5 text-xs text-cosmic-aurora">
            💡 空亡期间不宜开始新事务，适合完成已有工作
          </div>
        </>
      ) : (
        <div className="text-sm text-celestial-silver/60">
          距离下次空亡约 {voidInfo.duration.toFixed(1)} 小时
          <br />
          下一个星座: {voidInfo.nextSign}
        </div>
      )}
    </motion.div>
  );
}

// 行星时卡片
function PlanetaryHourCard({ hourInfo }: { hourInfo: PlanetaryHourInfo }) {
  const influenceColor =
    hourInfo.influence > 0
      ? 'text-green-400'
      : hourInfo.influence < 0
      ? 'text-red-400'
      : 'text-yellow-400';

  return (
    <motion.div
      className="p-4 rounded-lg bg-white/5"
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.1 }}
    >
      <div className="flex items-center justify-between mb-2">
        <div className="font-medium">⏰ 当前行星时</div>
        <div className={`text-sm font-bold ${influenceColor}`}>
          {hourInfo.influence > 0 ? '+' : ''}
          {hourInfo.influence.toFixed(1)}
        </div>
      </div>

      <div className="flex items-center gap-3 mb-3">
        <div className="text-3xl">{hourInfo.planetSymbol}</div>
        <div>
          <div className="text-lg font-bold">{hourInfo.planetName}时</div>
          <div className="text-sm text-celestial-silver/60">
            第 {hourInfo.planetaryHour} 个行星时
          </div>
        </div>
      </div>

      {hourInfo.bestFor && hourInfo.bestFor.length > 0 && (
        <div>
          <div className="text-xs text-celestial-silver/60 mb-1">适合活动</div>
          <div className="flex flex-wrap gap-1">
            {hourInfo.bestFor.map((activity, i) => (
              <span
                key={i}
                className="px-2 py-0.5 rounded bg-cosmic-nova/20 text-cosmic-nova text-xs"
              >
                {activity}
              </span>
            ))}
          </div>
        </div>
      )}
    </motion.div>
  );
}

// 月相图标
function MoonPhaseIcon({ phase }: { phase: string }) {
  const phaseIcons: Record<string, string> = {
    new: '🌑',
    waxing_crescent: '🌒',
    first_quarter: '🌓',
    waxing_gibbous: '🌔',
    full: '🌕',
    waning_gibbous: '🌖',
    last_quarter: '🌗',
    waning_crescent: '🌘',
  };

  return (
    <span className="text-3xl">{phaseIcons[phase] || '🌙'}</span>
  );
}

export default LunarInfo7934IJ;

