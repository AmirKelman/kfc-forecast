import { useQuery } from '@tanstack/react-query'
import { fetchForecasts } from '../api/queries'
import { useForecastStore } from '../store/forecastStore'
import DatePickerInput from './DatePickerInput'
import ForecastCard from './ForecastCard'

export default function ForecastPanel() {
  const { selectedStoreId, selectedDate } = useForecastStore()

  const { data, isLoading, isError, error } = useQuery({
    queryKey: ['forecasts', selectedStoreId, selectedDate],
    queryFn: () => fetchForecasts(selectedStoreId!, selectedDate),
    enabled: selectedStoreId !== null && !!selectedDate,
  })

  return (
    <main className="flex-1 overflow-y-auto p-6 min-w-0">
      {/* Top bar: store name + date picker */}
      <div className="flex flex-wrap items-center justify-between gap-4 mb-6">
        <div>
          {data?.store ? (
            <>
              <h2 className="text-2xl font-bold text-gray-900">{data.store.name}</h2>
              <p className="text-sm text-gray-400">{data.store.location}</p>
            </>
          ) : (
            <h2 className="text-2xl font-bold text-gray-300">
              {selectedStoreId ? '…' : 'Select a store'}
            </h2>
          )}
        </div>
        <DatePickerInput />
      </div>

      {/* States */}
      {!selectedStoreId && <NoStoreSelected />}

      {selectedStoreId && isLoading && <LoadingGrid />}

      {selectedStoreId && isError && (
        <ErrorState message={(error as Error)?.message ?? 'Unknown error'} />
      )}

      {selectedStoreId && data && data.forecasts.length === 0 && (
        <EmptyState date={selectedDate} />
      )}

      {selectedStoreId && data && data.forecasts.length > 0 && (
        <div className="grid grid-cols-1 xl:grid-cols-2 gap-4">
          {data.forecasts.map((f) => (
            <ForecastCard key={f.product.id} forecast={f} />
          ))}
        </div>
      )}
    </main>
  )
}

function NoStoreSelected() {
  return (
    <div className="flex flex-col items-center justify-center h-80 text-center">
      <div className="text-5xl mb-4">🏪</div>
      <p className="text-lg font-semibold text-gray-600">No store selected</p>
      <p className="text-sm text-gray-400 mt-1">Choose a store from the left sidebar</p>
    </div>
  )
}

function LoadingGrid() {
  return (
    <div className="grid grid-cols-1 xl:grid-cols-2 gap-4">
      {Array.from({ length: 6 }).map((_, i) => (
        <div key={i} className="bg-white rounded-xl border border-gray-200 p-4 h-52 animate-pulse">
          <div className="h-4 w-32 bg-gray-200 rounded mb-2" />
          <div className="h-3 w-16 bg-gray-100 rounded mb-4" />
          <div className="h-28 bg-gray-100 rounded" />
        </div>
      ))}
    </div>
  )
}

function ErrorState({ message }: { message: string }) {
  return (
    <div className="flex flex-col items-center justify-center h-80 text-center">
      <div className="text-5xl mb-4">⚠️</div>
      <p className="text-lg font-semibold text-red-600">Failed to load forecasts</p>
      <p className="text-sm text-gray-400 mt-1">{message}</p>
    </div>
  )
}

function EmptyState({ date }: { date: string }) {
  return (
    <div className="flex flex-col items-center justify-center h-80 text-center">
      <div className="text-5xl mb-4">📭</div>
      <p className="text-lg font-semibold text-gray-600">No forecasts for {date}</p>
      <p className="text-sm text-gray-400 mt-1">
        Click <strong>Generate Now</strong> to create predictions for tomorrow, or pick a
        different date.
      </p>
    </div>
  )
}
