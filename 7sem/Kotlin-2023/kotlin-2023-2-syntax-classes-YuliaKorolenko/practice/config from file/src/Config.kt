import java.io.InputStream
import kotlin.properties.ReadOnlyProperty
import kotlin.reflect.KProperty

class Config(configName: String) {
    private val resultMap = mutableMapOf<String, String>()

    init {
        val stream: InputStream = getResource(configName) ?: throw IllegalArgumentException("Can't open file")
        stream.reader().forEachLine { i ->
            if (!isCorrectString(i)) throw IllegalArgumentException("Not the only symbol")
            val splittedString = i.split("=")
            resultMap[splittedString[0].trim()] = splittedString[1].trim()
        }
    }

    fun isCorrectString(s: String): Boolean {
        return (s.length - s.replace("=", "").length) == 1
    }

    operator fun provideDelegate(
        receiver: Any?,
        property: KProperty<*>,
    ): ReadOnlyProperty<Any?, String> {
        if (!resultMap.containsKey(property.name)) {
            throw IllegalArgumentException("Don't have this value")
        }
        return ReadOnlyProperty { _, _ -> resultMap.getValue(property.name) }
    }
}
