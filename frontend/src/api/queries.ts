import client from './client'
import type { Store, ForecastResponse } from './types'

export const fetchStores = async (): Promise<Store[]> => {
  const { data } = await client.get<Store[]>('/stores')
  return data
}

export const fetchForecasts = async (
  storeId: number,
  date: string,
): Promise<ForecastResponse> => {
  const { data } = await client.get<ForecastResponse>('/forecasts', {
    params: { store_id: storeId, date },
  })
  return data
}

export const triggerGenerate = async (): Promise<void> => {
  await client.post('/forecasts/generate')
}
