package chatbot.dsl

import chatbot.api.Keyboard

class RowBuilder() {
    private var keyboardRow: MutableList<Keyboard.Button> = mutableListOf()

    operator fun String.unaryMinus() {
        keyboardRow.add(Keyboard.Button(text = this))
    }

    fun button(text: String) {
        keyboardRow.add(Keyboard.Button(text = text))
    }

    fun build(): MutableList<Keyboard.Button> {
        return keyboardRow
    }
}
