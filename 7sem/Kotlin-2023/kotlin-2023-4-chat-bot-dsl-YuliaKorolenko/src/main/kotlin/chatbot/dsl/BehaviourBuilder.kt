package chatbot.dsl

import chatbot.api.ChatContext
import chatbot.api.Message

@DslMarker
annotation class Behaviour

@Behaviour
class BehaviourBuilder() {
    val handlers = mutableListOf<HandlerBase>()

    fun onCommand(
        command: String,
        f: MessageProcessor<ChatContext?>,
    ) {
        handlers.add(HandlerBehaviour(f) { message -> message.text.startsWith("/$command") })
    }

    fun onMessage(f: MessageProcessor<ChatContext?>) {
        handlers.add(HandlerDefault(f))
    }

    fun onMessage(
        predicate: (Message) -> Boolean,
        f: MessageProcessor<ChatContext?>,
    ) {
        handlers.add(HandlerBehaviour(f, predicate))
    }

    fun onMessagePrefix(
        preffix: String,
        f: MessageProcessor<ChatContext?>,
    ) {
        handlers.add(HandlerBehaviour(f) { message -> message.text.startsWith(preffix) })
    }

    fun onMessageContains(
        text: String,
        f: MessageProcessor<ChatContext?>,
    ) {
        handlers.add(HandlerBehaviour(f) { message -> message.text.contains(text) })
    }

    fun onMessage(
        messageTextExactly: String,
        f: MessageProcessor<ChatContext?>,
    ) {
        handlers.add(HandlerBehaviour(f) { message -> message.text == messageTextExactly })
    }

    inline infix fun <reified T : ChatContext> T.into(configure: ContextBuilder<T>.() -> Unit) {
        val builder = ContextBuilder { if (it is T && this@into == it) it else null }
        configure(builder)
        handlers.addAll(builder.build())
    }

    inline fun <reified T : ChatContext> into(configure: ContextBuilder<T>.() -> Unit) {
        val builder = ContextBuilder { if (it is T) it else null }
        configure(builder)
        handlers.addAll(builder.build())
    }

    fun build(): List<HandlerBase> {
        return handlers
    }
}
