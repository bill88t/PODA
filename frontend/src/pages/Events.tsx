import { useEffect, useState } from "react";
import { useUser, type AppointmentView } from "../components/user/userContext.tsx"

export default function Events() {
    const usrCtx = useUser();
    if (usrCtx.user === null) throw Error("user is null in Events");
    const [events, setEvents] = useState<AppointmentView[]>([]);
    const [changed, setChanged] = useState<number>(0);
    const [create, setCreate] = useState<boolean>(false);


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
    }, [usrCtx, usrCtx.user, changed]);

    return (
        <>
            <title>Saloon PODA - Events</title>
            <h1>Events</h1>
            {events.map((e: AppointmentView) => (
                <div id={"apppontment" + e.appointment.id}>
                    <div>{e.fname},
                        {e.lname} -
                        {e.appointment.datetime.toLocaleString()},
                        <a onClick={
                            f => {
                                f.preventDefault();
                                usrCtx.deleteAppointment(e.userUuid, e.appointment.id);
                                setChanged(changed + 1 % 8);
                            }
                        }>Delete</a>
                    </div>
                </div>
            )
            )}
            <button className="span2"
                onClick={
                    f => {
                        f.preventDefault();
                        setCreate(!create);
                    }
                }
                id="create-appointment"
            >{create ?  "Abort" : "Create"}</button>
        </>
    )
}
