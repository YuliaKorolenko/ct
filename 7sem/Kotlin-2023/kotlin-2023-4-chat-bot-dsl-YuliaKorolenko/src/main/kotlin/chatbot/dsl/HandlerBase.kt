package chatbot.dsl

import chatbot.api.ChatContext
import chatbot.api.Message

sealed interface HandlerBase {
    fun checkerFun(message: Message, chatContext: ChatContext?): Boolean

    fun runHandler(proc: MessageProcessorContext<ChatContext?>)
}

open class HandlerDefault(var handlerBehaviour: MessageProcessor<ChatContext?>) : HandlerBase {
    override fun checkerFun(message: Message, chatContext: ChatContext?): Boolean = true

    override fun runHandler(proc: MessageProcessorContext<ChatContext?>) {
        handlerBehaviour(proc)
    }
}

class HandlerBehaviour(
    handlerBehaviour: MessageProcessor<ChatContext?>,
    var checker: ((Message) -> Boolean),
) : HandlerDefault(handlerBehaviour) {
    override fun checkerFun(message: Message, chatContext: ChatContext?): Boolean {
        return checker(message)
    }
}

class HandlerContext<T : ChatContext?>(
    var handlerBehaviour: MessageProcessor<T>,
    var checker: ((Message) -> Boolean),
    var cast: (ChatContext?) -> T?,
) : HandlerBase {
    override fun checkerFun(message: Message, chatContext: ChatContext?): Boolean {
        return checker(message) && cast(chatContext) != null
    }

    override fun runHandler(proc: MessageProcessorContext<ChatContext?>) {
        if (cast(proc.context) != null) {
            handlerBehaviour(
                MessageProcessorContext(
                    proc.message,
                    proc.client,
                    cast(proc.context)!!,
                    proc.setContext,
                ),
            )
        }
    }
}
