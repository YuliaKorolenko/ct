package airline.api

import kotlinx.datetime.Instant

sealed class Management(val flightId: String, val departureTime: Instant)

class ScheduleManagement(flightId: String, departureTime: Instant, val plane: Plane) :
    Management(flightId, departureTime)

class DelayManagement(flightId: String, departureTime: Instant, val actualDepartureTime: Instant) :
    Management(flightId, departureTime)

class CancelManagement(flightId: String, departureTime: Instant) : Management(flightId, departureTime)

class SetCheckManagement(flightId: String, departureTime: Instant, val checkInNumber: String) :
    Management(flightId, departureTime)

class SetGateManagement(flightId: String, departureTime: Instant, val gateNumber: String) :
    Management(flightId, departureTime)

class BuyManagement(
    flightId: String,
    departureTime: Instant,
    val seatNo: String,
    val passengerId: String,
    val passengerName: String,
    val passengerEmail: String,
) : Management(flightId, departureTime)
