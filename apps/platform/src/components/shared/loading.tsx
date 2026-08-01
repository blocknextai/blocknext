const Loading = () => {
  return (
    <div className="w-full h-full flex justify-center items-center">
      <span className="loading loading-infinity loading-xl"></span>
    </div>
  )
}
Loading.displayName = 'Loading'

const PageLoading = () => {
  return (
    <div className="flex-1 w-full flex justify-center items-center min-h-[calc(100vh-8rem)]">
      <span className="loading loading-infinity loading-xl"></span>
    </div>
  )
}
PageLoading.displayName = 'PageLoading'

export { Loading, PageLoading }
