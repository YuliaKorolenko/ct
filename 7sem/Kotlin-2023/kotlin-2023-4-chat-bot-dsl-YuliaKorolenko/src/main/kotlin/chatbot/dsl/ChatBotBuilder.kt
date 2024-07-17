package chatbot.dsl

import chatbot.api.*

@Behaviour
class ChatBotBuilder(val client: Client) {
    var logLevel: LogLevel = LogLevel.ERROR
    lateinit var handlers: List<HandlerBase>
    var contextManager: ChatContextsManager? = null

    fun use(info: LogLevel) {
        logLevel = info
    }

    fun use(manager: ChatContextsManager) {
        contextManager = manager
    }

    operator fun LogLevel.unaryPlus() {
        logLevel = this
    }

    fun behaviour(configure: BehaviourBuilder.() -> Unit) {
        val builder = BehaviourBuilder()
        configure(builder)
        this.handlers = builder.build()
    }

    fun build(): ChatBot {
        return object : ChatBot {
            override fun processMessages(message: Message) {
                for (handler in handlers) {
                    val context: ChatContext? = contextManager?.getContext(message.chatId)
                    if (handler.checkerFun(message, context)) {
                        handler.runHandler(
                            MessageProcessorContext(message, this@ChatBotBuilder.client, context) {
                                contextManager?.setContext(
                                    message.chatId,
                                    it,
                                )
                            },
                        )
                        break
                    }
                }
            }

            override val logLevel: LogLevel
                get() = this@ChatBotBuilder.logLevel
        }
    }
}

fun chatBot(
    client: Client,
    configure: ChatBotBuilder.() -> Unit,
): ChatBot {
    val builder = ChatBotBuilder(client)
    configure(builder)
    return builder.build()
}
