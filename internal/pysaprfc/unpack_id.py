import configparser

import pyrfc
import csv

from pyrfc import Connection

from pyrfc import ABAPApplicationError, ABAPRuntimeError, LogonError, CommunicationError, ExternalRuntimeError
import logging


# распаковка ЕО
def main():
    global output
    logger = logging.getLogger("unpack_id")
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
    config.read("sapnwrfc.cfg")
    config.sections()
    params_connection = config['connection']
    logger.info("Connecting to SAP RFC...")

    try:
        connection = pyrfc.Connection(**params_connection)
    #    result = connection.call('STFC_CONNECTION', REQUTEXT=u'Hello SAP!')
    #    print(result)

        # id_test_1000836.csv
        file_unpack_id = "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data/unpack_id.csv"

        # чтение списка выгруженных ЕО на редактирование
        rows = []
        with open(file_unpack_id, newline='') as file:
            csvreader = csv.reader(file, delimiter=',')
            header = next(csvreader)
            for row in csvreader:
                rows.append(row)
        # распаковка в цикле прочитанных ЕО
        for i in rows:
            output = connection.call('Z_IEXT_HU_UNPACKSNGLPOS', **{
                'UCODE': '21717',
                'PCODE': 'NEWPASSWORD1',
                'HUKEY': '0000000000' + i[0],  # '00000000000000015660',
                'ITEMUNPACK': {u'PACK_QTY': i[1]},

            })

        print(output)
        logger.info("Unpack ID: ", output)

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
    except (ABAPApplicationError, ABAPRuntimeError, ExternalRuntimeError):
        print("An error occurred.")
        logger.error("An error occurred.")
        logger.exception("Error!")
        raise


if __name__ == "__main__":
    main()
