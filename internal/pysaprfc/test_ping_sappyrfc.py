#!/usr/bin/env python3

import configparser
import pyrfc
from pyrfc import ABAPApplicationError, ABAPRuntimeError, LogonError, CommunicationError, ExternalRuntimeError
import logging


def main():
    logger = logging.getLogger("test_ping_sappyrfc")
    logger.setLevel(logging.INFO)

    # create the logging file handler
    fh = logging.FileHandler(
        "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/pyrfc_logging.log")
    formatter = logging.Formatter(
        '%(asctime)s - %(name)s - %(levelname)s - %(message)s')
    fh.setFormatter(formatter)

    # add handler to logger object
    logger.addHandler(fh)

    logger.info("Parsing cfg")
    config = configparser.ConfigParser()
    config.read(
        "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/sapnwrfc.cfg")
    config.sections()
    params_connection = config['connection']
    logger.info("Connecting to SAP RFC...")

    try:
        connection = pyrfc.Connection(**params_connection)
        result = connection.call('STFC_CONNECTION', REQUTEXT=u'Hello SAP!')
        print(result)
        logger.info(result)

        connection.close()

    except CommunicationError:
        print("Could not connect to server.")
        logger.error("Could not connect to server.")
        logger.exception("Error!")
        raise
    except LogonError:
        print("Could not log in. Wrong credentials?")
        logger.error("Could not log in. Wrong credentials?")
        logger.exception("Error!")
        raise
    except (ABAPApplicationError, ABAPRuntimeError):
        print("An error occurred.")
        logger.error("An error occurred.")
        logger.exception("Error!")
        raise


if __name__ == "__main__":
    main()
