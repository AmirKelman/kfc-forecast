import { useQuery } from '@tanstack/react-query'
import { fetchStores } from '../api/queries'
import { useForecastStore } from '../store/forecastStore'
import type { Store } from '../api/types'

export default function StoreList() {
  const { data: stores, isLoading, isError } = useQuery({
    queryKey: ['stores'],
    queryFn: fetchStores,
  })

  const { selectedStoreId, setSelectedStoreId } = useForecastStore()

  if (isLoading) {
    return (
      <aside className="w-64 shrink-0 bg-white border-r border-gray-200 p-4">
        <p className="text-xs font-semibold uppercase tracking-wider text-gray-400 mb-3">
          Stores
        </p>
        {[1, 2, 3].map((i) => (
          <div key={i} className="h-14 bg-gray-100 rounded-lg mb-2 animate-pulse" />
        ))}
      </aside>
    )
  }

  if (isError) {
    return (
      <aside className="w-64 shrink-0 bg-white border-r border-gray-200 p-4">
        <p className="text-sm text-red-500">Failed to load stores.</p>
      </aside>
    )
  }

  return (
    <aside className="w-64 shrink-0 bg-white border-r border-gray-200 p-4 overflow-y-auto">
      <p className="text-xs font-semibold uppercase tracking-wider text-gray-400 mb-3">
        Stores
      </p>
      <ul className="space-y-1">
        {stores?.map((store: Store) => {
          const isSelected = store.id === selectedStoreId
          return (
            <li key={store.id}>
              <button
                onClick={() => setSelectedStoreId(store.id)}
                className={`w-full text-left px-3 py-3 rounded-lg transition-all ${
                  isSelected
                    ? 'bg-red-50 border border-red-200 text-red-700'
                    : 'hover:bg-gray-50 text-gray-700 border border-transparent'
                }`}
              >
                <div className="flex items-center gap-2">
                  <span
                    className={`w-2 h-2 rounded-full shrink-0 ${
                      isSelected ? 'bg-red-500' : 'bg-gray-300'
                    }`}
                  />
                  <div className="min-w-0">
                    <p className="text-sm font-medium truncate">{store.name}</p>
                    <p className="text-xs text-gray-400 truncate">{store.location}</p>
                  </div>
                </div>
              </button>
            </li>
          )
        })}
      </ul>
    </aside>
  )
}
