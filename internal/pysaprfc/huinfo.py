import configparser
from operator import ne

import pyrfc
import csv

from pyrfc import ABAPApplicationError, ABAPRuntimeError, LogonError, CommunicationError, ExternalRuntimeError
import logging


def main():
    logger = logging.getLogger("huinfo")
    logger.setLevel(logging.INFO)

    global order

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
        logger.info("Connection to SAP RFC creating.")

        infoOrder = '/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test/test2_info_order.csv'
        id_file = "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test/test2_2_unpack_id.csv"
        # id_file = "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data/unpack_id.csv"
        sap_id_info = "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test/sap_id_info.csv"

        rowsinfoOrder = []
        with open(infoOrder, newline='') as file:
            csvreader = csv.DictReader(file, delimiter=',')
            for row in csvreader:
                rowsinfoOrder.append(row)
        for row in rowsinfoOrder:
            order = row['WO']

        rowsPana = []
        with open(id_file, newline='') as file:
            csvreader = csv.reader(file, delimiter=',')
            header = next(csvreader)
            for row in csvreader:
                rowsPana.append(row)
        arrID = []

        for i in rowsPana:
            print("i-", i[0])
            exid = '0000000000'+i[0]
            infoID = connection.call('Z_IEXT_HU_INFO', **{
                # 'UCODE': '21717',
                # 'PCODE': 'NEWPASSWORD1',
                'HUINFO': [
                    {
                        'EXIDV': exid,
                    }
                ]
            })
            id = get_info_id(infoID)

            if id is not None:
                arrID.append(id)
            else:
                print(f'error {exid}')
                logger.error(
                    f'При проверке получения данных ЕО нет информации о: ео - {exid}, заказ - {order}')
            # arrID.append(get_info_id(infoID))

        with open(sap_id_info, 'w', newline='') as wfile:
            idwriter = csv.writer(wfile, delimiter=',')
            idwriter.writerows(arrID)

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


def get_info_id(dict_value):
    for key, value in dict_value.items():
        if key == 'HUINFO':
            try:
                # + " " + value[0].get('MATNR', '')
                resID = value[0].get('EXIDV', '')
                resSAP = value[0].get('MATNR', '')
                resLOT = value[0].get('CHARG', '')
                resQty = value[0].get('VEMNG', '')
                resStock = value[0].get('LGORT', '')
                res = [resID, resSAP, resLOT, resQty, resStock]
            #   res = resID + "", "" + resSAP + "", "" + \
            #       resLOT + "", "" + str(resQty) + "", "" + resStock
                return res
            except IndexError:
                pass


if __name__ == "__main__":
    main()
