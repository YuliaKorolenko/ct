import kotlin.time.Duration
import kotlinx.coroutines.flow.*

fun Flow<Cutoff>.resultsFlow(): Flow<Results> {
    return scan(emptyMap<String, Duration>()) { acc, cutoff ->
        acc + (cutoff.number to cutoff.time)
    }
        .drop(1)
        .map { Results(it) }
}
