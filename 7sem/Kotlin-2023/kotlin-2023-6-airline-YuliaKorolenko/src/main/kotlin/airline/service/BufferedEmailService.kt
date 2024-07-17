package airline.service

import airline.api.Message
import kotlinx.coroutines.channels.Channel

class BufferedEmailService(private val emailService: EmailService) {
    private val messageChanel = Channel<Message>(100)

    suspend fun send(
        to: String,
        text: String,
    ) {
        messageChanel.send(Message(to, text))
    }

    suspend fun sendAllMessages() {
        while (true) {
            val message = messageChanel.receive()
            emailService.send(message.to, message.text)
        }
    }
}
