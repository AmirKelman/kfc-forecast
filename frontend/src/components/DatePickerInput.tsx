import { useForecastStore } from '../store/forecastStore'

export default function DatePickerInput() {
  const { selectedDate, setSelectedDate } = useForecastStore()

  return (
    <div className="flex items-center gap-2">
      <label
        htmlFor="forecast-date"
        className="text-sm font-medium text-gray-600 whitespace-nowrap"
      >
        Forecast date
      </label>
      <input
        id="forecast-date"
        type="date"
        value={selectedDate}
        onChange={(e) => e.target.value && setSelectedDate(e.target.value)}
        className="border border-gray-300 rounded-lg px-3 py-1.5 text-sm text-gray-800 focus:outline-none focus:ring-2 focus:ring-red-400 focus:border-transparent bg-white"
      />
    </div>
  )
}
