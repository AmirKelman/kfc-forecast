export interface Store {
  id: number
  name: string
  location: string
}

export interface Product {
  id: number
  name: string
  category: string
  price: number
}

export interface HourlyEntry {
  hour: number
  predicted_quantity: number
}

export interface ProductForecast {
  product: Product
  hourly: HourlyEntry[]
}

export interface ForecastResponse {
  store: Store
  date: string
  forecasts: ProductForecast[]
}
