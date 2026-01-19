import { useState, useEffect } from "react";
import { useUser, type AppointmentView } from "./user/userContext.tsx"

function CalendarEvent() {
    const ctxUsr = useUser();

    if (ctxUsr.user === null) throw Error("User for Calendar must exists");
    const [events, setEvents] = useState<AppointmentView[]>([]);

    useEffect(() => {
        (async () => {
            setEvents(
                await ctxUsr.viewAppointment(ctxUsr.user!.uuid, new Date())
            );
        })();
    }, [ctxUsr]);

    return (
        <div className="form">
            {events.map(e => (
                <>
                    <div className="field">
                        {e.fname} {e.lname}
                        <br />
                        {e.appointment.datetime.toLocaleString()}
                    </div>
                </>
            ))}
        </div>
    );
}

export function Calendar() {
    return <CalendarEvent />;
}
