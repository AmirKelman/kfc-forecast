import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { format, addDays } from 'date-fns'

interface ForecastStore {
  selectedStoreId: number | null
  selectedDate: string
  setSelectedStoreId: (id: number) => void
  setSelectedDate: (date: string) => void
}

const tomorrow = format(addDays(new Date(), 1), 'yyyy-MM-dd')

export const useForecastStore = create<ForecastStore>()(
  persist(
    (set) => ({
      selectedStoreId: null,
      selectedDate: tomorrow,
      setSelectedStoreId: (id) => set({ selectedStoreId: id }),
      setSelectedDate: (date) => set({ selectedDate: date }),
    }),
    { name: 'kfc-forecast-prefs' },
  ),
)
