type ChartPoint = {
  date: string
  avg: number
}

function formatDateLabel(date: string) {
  const parsed = new Date(date)
  if (Number.isNaN(parsed.getTime())) return date
  return parsed.toLocaleDateString('ja-JP', { month: '2-digit', day: '2-digit' })
}

export function PriceTrendChart({ data }: { data: ChartPoint[] }) {
  if (data.length === 0) {
    return (
      <div className="flex h-40 items-center justify-center rounded-xl border border-dashed border-slate-200 text-xs text-slate-400">
        📊 データなし
      </div>
    )
  }

  const width = 600
  const height = 180
  const paddingLeft = 42
  const paddingRight = 12
  const paddingTop = 12
  const paddingBottom = 28

  const avgPrices = data.map((item) => item.avg)
  const minValue = Math.min(...avgPrices)
  const maxValue = Math.max(...avgPrices)
  const range = maxValue - minValue || 1

  const scaleX = (index: number) =>
    paddingLeft + ((width - paddingLeft - paddingRight) * index) / Math.max(1, data.length - 1)
  const scaleY = (value: number) =>
    paddingTop + (height - paddingTop - paddingBottom) * (1 - (value - minValue) / range)

  // Grid lines (3 horizontal lines)
  const gridLines = [0, 0.5, 1].map((ratio) => {
    const value = minValue + range * ratio
    const y = scaleY(value)
    return { y, value }
  })

  // Line path
  const linePath =
    `M ${scaleX(0)},${scaleY(data[0].avg)} ` +
    data
      .slice(1)
      .map((item, index) => `L ${scaleX(index + 1)},${scaleY(item.avg)}`)
      .join(' ')

  // Area under line
  const areaPath =
    `M ${scaleX(0)},${height - paddingBottom} ` +
    `L ${scaleX(0)},${scaleY(data[0].avg)} ` +
    data
      .slice(1)
      .map((item, index) => `L ${scaleX(index + 1)},${scaleY(item.avg)}`)
      .join(' ') +
    ` L ${scaleX(data.length - 1)},${height - paddingBottom} Z`

  return (
    <svg viewBox={`0 0 ${width} ${height}`} className="w-full">
      <defs>
        <linearGradient id="areaGradient" x1="0" x2="0" y1="0" y2="1">
          <stop offset="0%" stopColor="#0ea5e9" stopOpacity="0.2" />
          <stop offset="100%" stopColor="#0ea5e9" stopOpacity="0.0" />
        </linearGradient>
      </defs>

      {/* Grid lines */}
      {gridLines.map((line, i) => (
        <g key={i}>
          <line
            x1={paddingLeft}
            y1={line.y}
            x2={width - paddingRight}
            y2={line.y}
            stroke="#e2e8f0"
            strokeWidth="1"
            opacity="0.5"
          />
          <text x={paddingLeft - 5} y={line.y + 3} fontSize="9" fill="#94a3b8" textAnchor="end">
            ¥{Math.round(line.value)}
          </text>
        </g>
      ))}

      {/* Area fill */}
      <path d={areaPath} fill="url(#areaGradient)" />

      {/* Main line */}
      <path
        d={linePath}
        fill="none"
        stroke="#0ea5e9"
        strokeWidth="2.5"
        strokeLinecap="round"
        strokeLinejoin="round"
      />

      {/* Data points */}
      {data.map((item, index) => (
        <g key={item.date} className="transition-all hover:scale-125 cursor-pointer">
          <circle cx={scaleX(index)} cy={scaleY(item.avg)} r="3.5" fill="#0ea5e9" />
          <circle cx={scaleX(index)} cy={scaleY(item.avg)} r="1.5" fill="#fff" />
        </g>
      ))}

      {/* Date labels */}
      {data.map((item, index) => {
        const shouldShow = data.length <= 7 || index % Math.ceil(data.length / 6) === 0 || index === data.length - 1
        if (shouldShow) {
          return (
            <text
              key={item.date}
              x={scaleX(index)}
              y={height - paddingBottom + 18}
              fontSize="9"
              fill="#64748b"
              textAnchor="middle"
            >
              {formatDateLabel(item.date)}
            </text>
          )
        }
        return null
      })}
    </svg>
  )
}
