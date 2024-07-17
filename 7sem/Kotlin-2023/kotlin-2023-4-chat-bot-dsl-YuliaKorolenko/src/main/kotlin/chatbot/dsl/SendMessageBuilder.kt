package chatbot.dsl

import chatbot.api.ChatId
import chatbot.api.Client
import chatbot.api.Keyboard
import chatbot.api.Keyboard.Markup
import chatbot.api.MessageId

class SendMessageBuilder(private var chatId: ChatId, private var client: Client) {
    var text: String = ""
    var replyTo: MessageId? = null
    private var keyboard: Keyboard? = null

    fun removeKeyboard() {
        keyboard = Keyboard.Remove
    }

    fun withKeyboard(configure: KeyboardBuilder.() -> Unit) {
        val builder = KeyboardBuilder()
        configure(builder)
        keyboard = builder.build()
    }

    fun build() {
        val keyMarkup = keyboard
        if (keyMarkup is Markup) {
            if (keyMarkup.keyboard.isEmpty() || keyMarkup.keyboard.all { it.isEmpty() }) {
                return
            }
        }
        if (text.isEmpty() && keyboard == null) {
            return
        }
        client.sendMessage(chatId, text, keyboard, replyTo)
    }
}
