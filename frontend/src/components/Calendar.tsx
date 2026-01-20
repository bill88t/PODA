import { useEffect, useState } from "react";
import { useUser, type AppointmentView } from "./user/userContext";

function CalendarEvent() {
    const usrCtx = useUser();
    const [events, setEvents] = useState<AppointmentView[]>([]);

    useEffect(() => {
        if (!usrCtx.user) return;

        const loadAppointments = async () => {
            const result = await usrCtx.viewAppointment(
                usrCtx.user!.uuid,
                new Date()
            );
            setEvents(result);
        };

        loadAppointments();
    }, [usrCtx, usrCtx.user]);

    if (!usrCtx.user) {
        return <div>No user loaded</div>;
    }

    return (
        <div className="form">
        {events.map(e => (
            <div className="field" key={e.appointment.id}>
                {e.fname} {e.lname}
                <br />
                {new Date(e.appointment.datetime).toLocaleString()}
            </div>
        ))}
        </div>
    );
}

export function Calendar() {
    return <CalendarEvent />;
}
