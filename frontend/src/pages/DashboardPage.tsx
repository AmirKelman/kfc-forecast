import Header from '../components/Header'
import StoreList from '../components/StoreList'
import ForecastPanel from '../components/ForecastPanel'

export default function DashboardPage() {
  return (
    <div className="flex flex-col h-screen bg-gray-50">
      <Header />
      <div className="flex flex-1 overflow-hidden">
        <StoreList />
        <ForecastPanel />
      </div>
    </div>
  )
}
