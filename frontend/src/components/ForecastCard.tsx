import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from 'recharts'
import type { ProductForecast } from '../api/types'

interface Props {
  forecast: ProductForecast
}

const CATEGORY_COLORS: Record<string, string> = {
  Burgers: '#E4002B',
  Buckets: '#FF6B35',
  Sides:   '#F5A623',
  Drinks:  '#4A90D9',
}

export default function ForecastCard({ forecast }: Props) {
  const { product, hourly } = forecast
  const color = CATEGORY_COLORS[product.category] ?? '#6B7280'

  const peak = Math.max(...hourly.map((h) => h.predicted_quantity))
  const total = hourly.reduce((s, h) => s + h.predicted_quantity, 0)

  return (
    <div className="bg-white rounded-xl border border-gray-200 p-4 shadow-sm hover:shadow-md transition-shadow">
      <div className="flex items-start justify-between mb-3">
        <div>
          <h3 className="text-sm font-semibold text-gray-900">{product.name}</h3>
          <span
            className="inline-block mt-0.5 text-xs px-2 py-0.5 rounded-full font-medium"
            style={{ backgroundColor: color + '18', color }}
          >
            {product.category}
          </span>
        </div>
        <div className="text-right">
          <p className="text-xs text-gray-400">Total predicted</p>
          <p className="text-lg font-bold text-gray-800">{Math.round(total)}</p>
        </div>
      </div>

      {hourly.length === 0 ? (
        <p className="text-xs text-gray-400 text-center py-6">No data for this hour range</p>
      ) : (
        <ResponsiveContainer width="100%" height={140}>
          <BarChart data={hourly} margin={{ top: 4, right: 4, left: -20, bottom: 0 }}>
            <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" vertical={false} />
            <XAxis
              dataKey="hour"
              tickFormatter={(h) => `${h}h`}
              tick={{ fontSize: 10, fill: '#9ca3af' }}
              axisLine={false}
              tickLine={false}
            />
            <YAxis
              tick={{ fontSize: 10, fill: '#9ca3af' }}
              axisLine={false}
              tickLine={false}
              domain={[0, Math.ceil(peak * 1.15)]}
            />
            <Tooltip
              formatter={(val) => [typeof val === 'number' ? Math.round(val) : val, 'Units']}
              labelFormatter={(h) => `${h}:00`}
              contentStyle={{
                fontSize: 12,
                borderRadius: 8,
                border: '1px solid #e5e7eb',
                boxShadow: '0 4px 6px -1px rgba(0,0,0,0.1)',
              }}
            />
            <Bar
              dataKey="predicted_quantity"
              fill={color}
              radius={[4, 4, 0, 0]}
              maxBarSize={24}
            />
          </BarChart>
        </ResponsiveContainer>
      )}
    </div>
  )
}
