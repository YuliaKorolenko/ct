import kotlin.coroutines.CoroutineContext
import kotlinx.coroutines.*

class ParallelEvaluator {
    suspend fun run(task: Task, n: Int, context: CoroutineContext) {
        try {
            coroutineScope {
                repeat(n) {
                    launch(context) {
                        task.run(it)
                    }
                }
            }
        } catch (e: Exception) {
            throw TaskEvaluationException(e)
        }
    }
}
