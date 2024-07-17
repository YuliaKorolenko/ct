package airline

import airline.api.*
import airline.service.*
import kotlin.time.Duration.Companion.minutes
import kotlinx.coroutines.*
import kotlinx.coroutines.flow.*
import kotlinx.datetime.Clock
import kotlinx.datetime.Instant

class AirlineApplication(private val config: AirlineConfig, private val emailService: EmailService) {
    val bookingService: BookingService =
        object : BookingService {
            override val flightSchedule: List<FlightInfo>
                get() =
                    flights.value.map { flight ->
                        FlightInfo(
                            flightId = flight.flightId,
                            departureTime = flight.departureTime,
                            isCancelled = flight.isCancelled,
                            actualDepartureTime = flight.actualDepartureTime,
                            checkInNumber = flight.checkInNumber,
                            gateNumber = flight.gateNumber,
                            plane = flight.plane,
                        )
                    }.filter { flight ->
                        (flight.actualDepartureTime - config.ticketSaleEndTime) > Clock.System.now()
                    }

            override fun freeSeats(
                flightId: String,
                departureTime: Instant,
            ): Set<String> {
                val curFlight =
                    flights.value.find { flight ->
                        flight.flightId == flightId && flight.departureTime == departureTime &&
                            !flight.isCancelled && flight.departureTime > Clock.System.now()
                    }
                return curFlight?.let { f ->
                    val curSeats = f.plane.seats
                    curSeats.filter { seatNo -> !f.tickets.containsKey(seatNo) }.toSet()
                } ?: emptySet()
            }

            override suspend fun buyTicket(
                flightId: String,
                departureTime: Instant,
                seatNo: String,
                passengerId: String,
                passengerName: String,
                passengerEmail: String,
            ) {
                serviceFlow.emit(
                    BuyManagement(flightId, departureTime, seatNo, passengerId, passengerName, passengerEmail),
                )
            }

        }

    val managementService: AirlineManagementService =
        object : AirlineManagementService {
            override suspend fun scheduleFlight(
                flightId: String,
                departureTime: Instant,
                plane: Plane,
            ) {
                serviceFlow.emit(ScheduleManagement(flightId, departureTime, plane))
            }

            override suspend fun delayFlight(
                flightId: String,
                departureTime: Instant,
                actualDepartureTime: Instant,
            ) {
                serviceFlow.emit(DelayManagement(flightId, departureTime, actualDepartureTime))
            }

            override suspend fun cancelFlight(
                flightId: String,
                departureTime: Instant,
            ) {
                serviceFlow.emit(CancelManagement(flightId, departureTime))
            }

            override suspend fun setCheckInNumber(
                flightId: String,
                departureTime: Instant,
                checkInNumber: String,
            ) {
                serviceFlow.emit(SetCheckManagement(flightId, departureTime, checkInNumber))
            }

            override suspend fun setGateNumber(
                flightId: String,
                departureTime: Instant,
                gateNumber: String,
            ) {
                serviceFlow.emit(SetGateManagement(flightId, departureTime, gateNumber))
            }

        }

    private val bufferedService = BufferedEmailService(emailService)

    var flights: MutableStateFlow<List<Flight>> = MutableStateFlow(emptyList())

    private val passengerNotificationService = PassengerNotificationService(bufferedService, flights)

    @OptIn(FlowPreview::class)
    fun airportInformationDisplay(coroutineScope: CoroutineScope): StateFlow<InformationDisplay> {
        return flights.map { flightList ->
            val infos =
                flightList.map { flight: Flight ->
                    FlightInfo(
                        flight.flightId,
                        flight.departureTime,
                        flight.isCancelled,
                        flight.actualDepartureTime,
                        flight.checkInNumber,
                        flight.gateNumber,
                        flight.plane,
                    )
                }
            InformationDisplay(infos)
        }.sample(config.displayUpdateInterval)
            .stateIn(coroutineScope, SharingStarted.Eagerly, InformationDisplay(emptyList()))
    }

    val airportAudioAlerts: Flow<AudioAlerts> =
        flow {
            while (true) {
                flights.value.filter { !it.isCancelled && (it.gateNumber != null || it.checkInNumber != null) }
                    .forEach { flight ->
                        val curTime = Clock.System.now()
                        val registrationOpeningTime = flight.actualDepartureTime - config.registrationOpeningTime
                        val registrationClosingTime = flight.actualDepartureTime - config.registrationClosingTime
                        val boardingOpeningTime = flight.actualDepartureTime - config.boardingOpeningTime
                        val boardingClosingTime = flight.actualDepartureTime - config.boardingClosingTime

                        if (flight.gateNumber != null) {
                            if (curTime in (boardingOpeningTime..boardingOpeningTime + 3.minutes)) {
                                emit(AudioAlerts.BoardingOpened(flight.flightId, flight.gateNumber))
                            }
                            if (curTime in (boardingClosingTime - 3.minutes..boardingClosingTime)) {
                                emit(AudioAlerts.BoardingClosing(flight.flightId, flight.gateNumber))
                            }
                        }

                        if (flight.checkInNumber != null) {
                            if (curTime in (registrationOpeningTime..registrationOpeningTime + 3.minutes)) {
                                emit(AudioAlerts.RegistrationOpen(flight.flightId, flight.checkInNumber))
                            }
                            if (curTime in (registrationClosingTime - 3.minutes..registrationClosingTime)) {
                                emit(AudioAlerts.RegistrationClosing(flight.flightId, flight.checkInNumber))
                            }
                        }
                    }

                delay(config.audioAlertsInterval)
            }
        }

    private val serviceFlow = MutableSharedFlow<Management>(extraBufferCapacity = 1000000)

    private suspend fun changeFlight(
        updatedFlights: MutableList<Flight>,
        existingFlight: Flight?,
        newFlight: Flight?,
        senNot: (Flight) -> NotificationMessage,
    ) {
        if (existingFlight != null && newFlight != null) {
            updatedFlights.remove(existingFlight)
            updatedFlights.add(newFlight)
            flights.value = listOf(*updatedFlights.toTypedArray())
            passengerNotificationService.addNotification(senNot(newFlight))
        }
    }

    suspend fun run() {
        coroutineScope {
            launch { bufferedService.sendAllMessages() }
            launch { passengerNotificationService.sendMessages() }
            launch {
                serviceFlow.collect { sAction: Management ->
                    val existingFlight =
                        flights.value.find { flight ->
                            flight.flightId == sAction.flightId &&
                                flight.departureTime == sAction.departureTime &&
                                !flight.isCancelled
                        }
                    val updatedFlights = flights.value.toMutableList()
                    when (sAction) {
                        is ScheduleManagement -> {
                            if (existingFlight == null) {
                                val newFlight =
                                    Flight(
                                        sAction.flightId,
                                        sAction.departureTime,
                                        false,
                                        sAction.departureTime,
                                        null,
                                        null,
                                        sAction.plane,
                                    )
                                updatedFlights.add(newFlight)
                                flights.value = updatedFlights.toList()
                            }
                        }

                        is DelayManagement -> {
                            val newFlight = existingFlight?.copy(actualDepartureTime = sAction.actualDepartureTime)
                            changeFlight(updatedFlights, existingFlight, newFlight) { i -> DelayNotification(i) }
                        }

                        is CancelManagement -> {
                            val newFlight = existingFlight?.copy(isCancelled = true)
                            changeFlight(updatedFlights, existingFlight, newFlight) { i -> CancelNotification(i) }
                        }

                        is SetCheckManagement -> {
                            val newFlight = existingFlight?.copy(checkInNumber = sAction.checkInNumber)
                            changeFlight(updatedFlights, existingFlight, newFlight) { i -> CheckInNotification(i) }
                        }

                        is SetGateManagement -> {
                            val newFlight = existingFlight?.copy(gateNumber = sAction.gateNumber)
                            changeFlight(updatedFlights, existingFlight, newFlight) { i -> SetGateNotification(i) }
                        }

                        is BuyManagement -> {
                            var message =
                                "Sorry, you couldn't buy a ticket ${sAction.seatNo} " +
                                    "for flight ${sAction.flightId}. There's been a mistake."
                            if (existingFlight != null && sAction.seatNo in existingFlight.plane.seats &&
                                sAction.seatNo !in existingFlight.tickets &&
                                existingFlight.actualDepartureTime - config.ticketSaleEndTime > Clock.System.now() &&
                                !existingFlight.tickets.values.any {
                                        ticket ->
                                    ticket.passengerId == sAction.passengerId
                                }
                            ) {
                                existingFlight.tickets[sAction.seatNo] =
                                    Ticket(
                                        sAction.flightId,
                                        sAction.departureTime,
                                        sAction.seatNo,
                                        sAction.passengerId,
                                        sAction.passengerName,
                                        sAction.passengerEmail,
                                    )
                                message = "You successfully bought ticket for flight ${sAction.flightId}," +
                                    " your seat number ${sAction.seatNo}."
                            }
                            bufferedService.send(
                                sAction.passengerEmail,
                                "Dear, ${sAction.passengerName}. " + message,
                            )
                        }
                    }
                }
            }
        }
    }
}
