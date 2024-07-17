import kotlinx.coroutines.*
import kotlinx.coroutines.withContext

class SequentialProcessor(private val handler: (String) -> String) : TaskProcessor {
    @OptIn(ExperimentalCoroutinesApi::class, DelicateCoroutinesApi::class)
    private val singleThreadContext = newSingleThreadContext("SequentialProcessorThreadContext")
    override suspend fun process(argument: String): String {
        return withContext(singleThreadContext) {
            handler(argument)
        }
    }
}
