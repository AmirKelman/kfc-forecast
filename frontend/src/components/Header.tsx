import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { triggerGenerate } from '../api/queries'

export default function Header() {
  const queryClient = useQueryClient()
  const [showSuccess, setShowSuccess] = useState(false)

  const { mutate, isPending } = useMutation({
    mutationFn: triggerGenerate,
    onSuccess: () => {
      setShowSuccess(true)
      setTimeout(() => {
        queryClient.invalidateQueries({ queryKey: ['forecasts'] })
        setShowSuccess(false)
      }, 3000)
    },
  })

  return (
    <header className="bg-white border-b border-gray-200 px-6 py-4 flex items-center justify-between shadow-sm">
      <div className="flex items-center gap-3">
        <div className="w-10 h-10 bg-red-600 rounded-xl flex items-center justify-center text-xl shadow-sm">
          🍗
        </div>
        <div>
          <h1 className="text-xl font-bold text-gray-900 leading-tight">
            KFC Sales Forecast
          </h1>
          <p className="text-xs text-gray-400">Daily predictions · per store & product</p>
        </div>
      </div>

      <div className="flex items-center gap-3">
        {showSuccess && (
          <span className="text-sm text-green-600 font-medium animate-pulse">
            ✓ Generating for tomorrow…
          </span>
        )}
        <button
          onClick={() => mutate()}
          disabled={isPending}
          className="flex items-center gap-2 px-4 py-2 bg-red-600 hover:bg-red-700 disabled:bg-red-300 text-white text-sm font-semibold rounded-lg transition-colors shadow-sm"
        >
          {isPending ? (
            <>
              <span className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
              Generating…
            </>
          ) : (
            <>▶ Generate Now</>
          )}
        </button>
      </div>
    </header>
  )
}
