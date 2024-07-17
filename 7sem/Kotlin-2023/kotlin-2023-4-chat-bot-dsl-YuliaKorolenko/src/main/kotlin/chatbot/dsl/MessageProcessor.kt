package chatbot.dsl

import chatbot.api.*

@Behaviour
class MessageProcessorContext<C : ChatContext?>(
    val message: Message,
    val client: Client,
    val context: C,
    val setContext: (c: ChatContext?) -> Unit,
) {
    fun sendMessage(chatId: ChatId, configure: SendMessageBuilder.() -> Unit) {
        val builder = SendMessageBuilder(chatId, client)
        configure(builder)
        return builder.build()
    }
}

typealias MessageProcessor<C> = MessageProcessorContext<C>.() -> Unit
