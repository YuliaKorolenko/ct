package chatbot.dsl

import chatbot.api.ChatContext
import chatbot.api.Message

@Behaviour
class ContextBuilder<T : ChatContext?>(private val contextChecker: (ChatContext?) -> T?) {
    private val handlers = mutableListOf<HandlerContext<T>>()

    fun onCommand(
        command: String,
        f: MessageProcessor<T>,
    ) {
        handlers.add(HandlerContext(f, { message -> message.text.startsWith("/$command") }, contextChecker))
    }

    fun onMessage(f: MessageProcessor<T>) {
        handlers.add(HandlerContext(f, { true }, contextChecker))
    }

    fun onMessage(
        predicate: (Message) -> Boolean,
        f: MessageProcessor<T>,
    ) {
        handlers.add(HandlerContext(f, predicate, contextChecker))
    }

    fun onMessagePrefix(
        preffix: String,
        f: MessageProcessor<T>,
    ) {
        handlers.add(HandlerContext(f, { message -> message.text.startsWith(preffix) }, contextChecker))
    }

    fun onMessageContains(
        text: String,
        f: MessageProcessor<T>,
    ) {
        handlers.add(HandlerContext(f, { message -> message.text.contains(text) }, contextChecker))
    }

    fun onMessage(
        messageTextExactly: String,
        f: MessageProcessor<T>,
    ) {
        handlers.add(HandlerContext(f, { message -> message.text == messageTextExactly }, contextChecker))
    }

    fun build(): List<HandlerContext<out ChatContext?>> {
        return this.handlers
    }
}
