package chatbot.dsl

import chatbot.api.*

class KeyboardBuilder() {
    var oneTime: Boolean = false
    var keyboard: MutableList<MutableList<Keyboard.Button>> = mutableListOf()

    fun row(configure: RowBuilder.() -> Unit) {
        val builder = RowBuilder()
        configure(builder)
        keyboard.add(builder.build())
    }

    fun build(): Keyboard.Markup {
        return Keyboard.Markup(oneTime, keyboard)
    }
}
