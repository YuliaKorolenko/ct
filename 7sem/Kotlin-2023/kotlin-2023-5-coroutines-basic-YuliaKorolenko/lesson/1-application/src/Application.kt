import kotlinx.coroutines.*

fun CoroutineScope.runApplication(
    runUI: suspend () -> Unit,
    runApi: suspend () -> Unit,
) {
    val jobApi = launch {
        var isReady = false
        while (!isReady) {
            try {
                runApi()
                isReady = true
            } catch (e: Exception) {
                delay(1000)
            }
        }
    }

    launch {
        try {
            runUI()
        } catch (e: Exception) {
            jobApi.cancel()
            throw e
        }
    }
}
