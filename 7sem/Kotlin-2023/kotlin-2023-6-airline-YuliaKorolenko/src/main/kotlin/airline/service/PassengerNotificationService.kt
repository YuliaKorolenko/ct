package airline.service

import airline.api.Flight
import kotlinx.coroutines.channels.Channel
import kotlinx.coroutines.flow.MutableStateFlow

class PassengerNotificationService(
    private val bufService: BufferedEmailService,
    private val flights: MutableStateFlow<List<Flight>>,
) {
    private val notificationChanel = Channel<NotificationMessage>(100)

    suspend fun sendMessages() {
        while (true) {
            val notification = notificationChanel.receive()
            val curFlight =
                flights.value.find {
                    it.departureTime == notification.flight.departureTime &&
                        it.flightId == notification.flight.flightId
                } ?: return
            val notPassenger = curFlight.tickets.values
            for (person in notPassenger) {
                bufService.send(
                    person.passengerEmail,
                    "Dear, ${person.passengerName}. " + notification.message,
                )
            }
        }
    }

    suspend fun addNotification(notification: NotificationMessage) {
        notificationChanel.send(notification)
    }
}

sealed class NotificationMessage(val flight: Flight) {
    abstract val message: String
}

class DelayNotification(flight: Flight) : NotificationMessage(flight) {
    override val message: String =
        "Your flight with number ${flight.flightId} delayed. The previous time of sending: ${flight.departureTime}. " +
            "New departure time ${flight.actualDepartureTime}."
}

class CheckInNotification(flight: Flight) : NotificationMessage(flight) {
    override val message: String =
        "Your flight with number ${flight.flightId}. " +
            "The flight received its check-in number: ${flight.checkInNumber}." +
            " We remind you that your plane departs at ${flight.actualDepartureTime}"
}

class CancelNotification(flight: Flight) : NotificationMessage(flight) {
    override val message: String =
        "Your flight with number ${flight.flightId} has cancelled. " +
            "The money for the ticket will be refunded within 5 days"
}

class SetGateNotification(flight: Flight) : NotificationMessage(flight) {
    override val message: String =
        "Your flight with number ${flight.flightId} will depart from gate ${flight.gateNumber}"
}
