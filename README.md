# GoInit

**GoInit** è un init system scritto in Go. Fornisce un gestore di servizi modulare e facilmente estendibile.

## ✨ Caratteristiche

- 📦 Definizione servizi in JSON
- 🔄 Gestione dei servizi tramite FIFO client-server
- 📁 Logging modulare con LMD (Log Manager Daemon)
- 🔧 Architettura componibile (SMNG, LMD, CMCN)
- 🧪 Supporto per attivazione temporizzata (`At`)
- 🚧 Futuro supporto a DBus e segnali POSIX (`SIGTERM`, ecc.)

## 🧩 Architettura

GoInit è composto da tre componenti principali:

- **SERVICES** (Service Manager): carica, attiva, disattiva e monitora i servizi.
- **LMD** (Log Manager Daemon): gestisce i log in formato leggibile e colorato (es. `[ ERROR ]` in rosso).

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
- `User`: con quale utente eseguire il processo
- `At`: quando avviarlo (`0` = All'avvio)

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
Se vuoi contribuire, apri un’[issue](https://github.com/Hacker8P/GoInit/issues) o una PR.

---

**Licenza**: GPL 3.0
**Autore**: Hacker8P