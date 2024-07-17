sealed class IntResult {
    data class Ok(val value: Int) : IntResult()

    data class Error(val reason: String) : IntResult()

    fun getOrDefault(i: Int): Int {
        return when (this) {
            is Ok -> this.value
            is Error -> i
        }
    }

    fun getOrNull(): Int? {
        return when (this) {
            is Ok -> this.value
            is Error -> null
        }
    }

    fun getStrict(): Int {
        return when (this) {
            is Ok -> this.value
            is Error -> throw NoResultProvided("custom reason")
        }
    }
}

fun safeRun(function: () -> Int): IntResult {
    return try {
        IntResult.Ok(function())
    } catch (e: Exception) {
        IntResult.Error(e.message ?: "undefined")
    }
}

class NoResultProvided(message: String) : NoSuchElementException(message)
