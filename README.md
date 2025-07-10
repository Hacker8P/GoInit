# GoInit

**GoInit** è un semplice init system scritto in Go. Fornisce un gestore di servizi modulare, leggero e facilmente estendibile, pensato per ambienti containerizzati, sistemi embedded, o come base per init personalizzati.

## ✨ Caratteristiche

- 📦 Definizione servizi in JSON
- 🔄 Gestione dei servizi tramite FIFO client-server
- 📁 Logging modulare con LMD (Log Manager Daemon)
- 🔧 Architettura componibile (SMNG, LMD, CMCN)
- 🧪 Supporto per attivazione temporizzata (`At`)
- 🚧 Futuro supporto a DBus e segnali POSIX (`SIGTERM`, ecc.)

## 🧩 Architettura

GoInit è composto da tre componenti principali:

- **SMNG** (Service Manager): carica, attiva, disattiva e monitora i servizi.
- **LMD** (Log Manager Daemon): gestisce i log in formato leggibile e colorato (es. `[ ERROR ]` in rosso).
- **CMCN** (Communication Manager - *DEPRECATO*): vecchio gestore della comunicazione, verrà sostituito da un modulo DBus o altro.

## ⚙️ Formato Servizi

Ogni servizio è descritto in formato JSON:

```json
{
  "Name": "Echo",
  "Command": "/usr/bin/echo Ciao",
  "Active": true,
  "At": 0
}
```

- `Name`: nome univoco del servizio
- `Command`: comando da eseguire
- `Active`: se avviarlo automaticamente all'avvio
- `At`: tempo in secondi dopo il boot per avviarlo (`0` = subito)

## 🔌 Comunicazione

GoInit utilizza una FIFO per la comunicazione tra client e demone. I messaggi sono in formato JSON.

Esempio di richiesta via client:

```json
{
  "action": "start",
  "service": "Echo"
}
```

Per inviare comandi:

```bash
echo '{"action":"status"}' > /tmp/goinit.fifo
```

## 📄 Logging

Il modulo **LMD** stampa i log in due stream distinti:

- `stdout`: log informativi
- `stderr`: errori con evidenziazione `[ ERROR ]` (es. rosso)

## 🚀 Avvio

Per compilare e lanciare GoInit:

```bash
go build -o goinit
sudo ./goinit
```

## ✅ To Do

- [ ] Gestione segnali POSIX (`SIGTERM`, `SIGINT`)
- [ ] Reload dinamico dei servizi
- [ ] Interfaccia DBus
- [ ] Monitoraggio avanzato dello stato
- [ ] File di configurazione esterni
- [ ] Test unitari e di integrazione

## 📦 Requisiti

- Go 1.20 o superiore
- Linux (richiesta compatibilità con FIFO e syscall base)

## 🤝 Contribuire

Pull request, issue e suggerimenti sono benvenuti!  
Se vuoi contribuire, apri un’[issue](https://github.com/tuo-utente/goinit/issues) o una PR.

---

**Licenza**: MIT  
**Autore**: [Il tuo nome o nickname]
