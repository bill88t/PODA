# PODA
Πολυτεχνικής Σχολής Παπουτσοπωλείο Δυτικής Αττικής

## Build on bread right now

```bash
sudo make breadrun
```

Η από πάνω εντολή φτιάχνει και τρέχει τον Server στην πόρτα 80

```bash
make run
```

Η από πάνω εντολή τρέχει στην πόρτα 5137

```
make clean
```

Η από πάνω εντολή καθαρίζει τα εκτελέσιμα

```bash
make build
```

Δημιουργείτε το main πρόγραμμα και χτίζει το frontend στον φάκελο
`frontend/dist`
όπου και γίνεται αναμετάδοση από το fiber σαν static asset

Το index.html είναι ορισμένο ως το wildcard ώστε να μπορεί να
αναπαράγει τα path του
frontend


