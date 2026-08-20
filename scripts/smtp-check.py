#!/usr/bin/env python3
"""Send one test email using the SMTP settings in .env, and say plainly what happened.

Run this before restarting the stack. It talks to the relay directly, so a bad
password or a blocked port fails here in two seconds instead of turning into a
silent "the candidate never got their link" during a campaign.

    ./scripts/smtp-check.py you@gmail.com
"""
import os
import smtplib
import ssl
import sys
from email.message import EmailMessage
from pathlib import Path

GREEN, RED, YELLOW, DIM, OFF = "\033[32m", "\033[31m", "\033[33m", "\033[2m", "\033[0m"
ok = lambda m: print(f"  {GREEN}OK{OFF}  {m}")
bad = lambda m: print(f"  {RED}!!{OFF}  {m}")
note = lambda m: print(f"      {DIM}{m}{OFF}")


def load_env(path: Path) -> dict:
    """Minimal .env reader — the same file docker compose reads."""
    values = {}
    if not path.exists():
        return values
    for line in path.read_text().splitlines():
        line = line.strip()
        if not line or line.startswith("#") or "=" not in line:
            continue
        key, _, val = line.partition("=")
        values[key.strip()] = val.strip().strip('"').strip("'")
    return values


def main() -> int:
    if len(sys.argv) < 2:
        print(__doc__)
        return 2
    to = sys.argv[1]

    env = {**load_env(Path(__file__).resolve().parent.parent / ".env"), **os.environ}
    host = env.get("SMTP_HOST", "")
    port = int(env.get("SMTP_PORT") or 587)
    user = env.get("SMTP_USER", "")
    password = env.get("SMTP_PASS", "")
    sender = env.get("SMTP_FROM") or user

    print(f"\n\033[1mSMTP check\033[0m  {host}:{port}  as {user or '(no auth)'}  →  {to}\n")

    if not host:
        bad("SMTP_HOST is not set, so the stack will only log links, never send them.")
        return 1
    if host == "mailpit":
        bad("SMTP_HOST is still 'mailpit' — mail is caught locally and never reaches a real inbox.")
        note("Set the real values in .env to send to a genuine address.")
        return 1
    if not sender:
        bad("SMTP_FROM (or SMTP_USER) is empty — there is no From address to send as.")
        return 1

    # Gmail rewrites the From header to the authenticated account, so a mismatch
    # silently changes who the candidate sees the mail from.
    if "gmail.com" in host and user and sender.lower() != user.lower():
        print(f"  {YELLOW}??{OFF}  SMTP_FROM ({sender}) differs from SMTP_USER ({user}).")
        note("Gmail rewrites From to the authenticated account unless it is a verified alias.")

    msg = EmailMessage()
    msg["From"] = sender
    msg["To"] = to
    msg["Subject"] = "Knovate SMTP check"
    msg.set_content(
        "If you are reading this, the scholarship funnel can email candidates.\n\n"
        f"Sent through {host}:{port} as {user or 'an unauthenticated relay'}.\n"
    )

    try:
        if port == 465:
            server = smtplib.SMTP_SSL(host, port, context=ssl.create_default_context(), timeout=20)
        else:
            server = smtplib.SMTP(host, port, timeout=20)
            server.ehlo()
            if server.has_extn("starttls"):
                server.starttls(context=ssl.create_default_context())
                server.ehlo()
                ok("STARTTLS negotiated")
            else:
                print(f"  {YELLOW}??{OFF}  the relay does not offer STARTTLS — traffic is in the clear")
        with server:
            ok(f"connected to {host}:{port}")
            if user:
                server.login(user, password)
                ok("authenticated")
            else:
                note("no SMTP_USER set — sending without authentication")
            server.send_message(msg)
            ok(f"accepted for delivery to {to}")
    except smtplib.SMTPAuthenticationError as e:
        bad(f"the relay rejected those credentials: {e.smtp_error.decode(errors='replace')[:160]}")
        if "gmail" in host:
            note("Gmail needs an App Password, not the account password.")
            note("Turn on 2-Step Verification, then Google Account -> Security -> App passwords.")
            note("Paste the 16 characters into SMTP_PASS (spaces are fine).")
        return 1
    except (smtplib.SMTPConnectError, OSError) as e:
        bad(f"could not reach {host}:{port} — {e}")
        note("Usually a firewall or an ISP blocking outbound 587. Try port 465.")
        return 1
    except smtplib.SMTPException as e:
        bad(f"the relay refused the message: {e}")
        return 1

    print(f"\n{GREEN}Delivered.{OFF} Check {to} — including spam, which is where a new sender usually lands.\n")
    return 0


if __name__ == "__main__":
    sys.exit(main())
