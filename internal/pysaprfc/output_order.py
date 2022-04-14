import configparser
import csv

import pyrfc

from pyrfc import Connection

from pyrfc import ABAPApplicationError, ABAPRuntimeError, LogonError, CommunicationError
import logging


# выпуск изделия
def main():
    # log_file = 'logging.log'
    # fl = open(log_file, 'a+')
    # fl.close()
    # logging.basicConfig(filename="logging.log", level=logging.INFO)
    global paramsATHDRLEVELS, paramsGOODSMOVEMENTS, SAP_ORDER
    logger = logging.getLogger("output_order")
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
        "/ home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/sapnwrfc.cfg")
    config.sections()
    params_connection = config['connection']
    logger.info("Connecting to SAP RFC...")

    try:

        # while True:
        connection = pyrfc.Connection(**params_connection)
       # result = connection.call('STFC_CONNECTION', REQUTEXT=u'Hello SAP!')
       # print(result)
        logger.info("Connection to SAP RFC creating.")
       # resultTime = connection.call(
       #     'WEEK_GET_FIRST_DAY', **{'WEEK': '201825'})
       # print(resultTime)

        infoOrder = '/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data/info_order.csv'
        # 'info_material_order.csv'
        infomaterialOrder = '/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data/wo_component.csv'

        rowsinfoOrder = []
        with open(infoOrder, newline='') as file:
            csvreader = csv.DictReader(file, delimiter=',')
            for row in csvreader:
                rowsinfoOrder.append(row)
        for row in rowsinfoOrder:
            paramsATHDRLEVELS = [

                {
                    # '000001000825s',  # 000001000825
                    'ORDERID': '00000' + row['WO'],
                    'YIELD': row['Qty'],  # '1.00',
                    'POSTG_DATE': row['Date'],  # '20220211',
                    'FIN_CONF': '',
                    'CLEAR_RES': '',
                }
            ]
        print("rowsinfoOrder:", type(rowsinfoOrder))
        sapORDER = [sub['WO'] for sub in rowsinfoOrder]
        print(str(sapORDER))
        for i in sapORDER:
            print("i", i)
            SAP_ORDER = i
        # print("sap order", rowsinfoOrder['WO'])

        infoMaterialOrder = []
        with open(infomaterialOrder, newline='') as fileMaterial:
            csvreader = csv.DictReader(fileMaterial, delimiter=',')
            for row in csvreader:
                infoMaterialOrder.append(row)

        order_info = connection.call('Z_IEXT_PRODORD_INFO', **{
            'AUFNR': '00000' + SAP_ORDER,
            # '000001000836', # str('00000' + str(sapORDER)),  # '00000' + row['WO'],  # orderSAP,  # 000001000825
            'UCODE': '21717',
            'PCODE': 'NEWPASSWORD1',
        }
        )
        productsap = order_info['PRODUCT']
        print("productsap", productsap)
        sap_order = order_info['RESITEMS']
        print("sap_order", sap_order)
        # breakpoint()
        paramsGOODSMOVEMENTS = []
        for row in rowsinfoOrder:
            paramsGOODSMOVEMENTS.append({
                'MATERIAL': productsap,  # '000000000003100302',
                'PLANT': 'SL00',
                'MOVE_TYPE': '131',
                'ENTRY_QNT': row['Qty'],  # '1',
                'ENTRY_UOM': 'ST',
                'ORDERID': '00000' + row['WO'],  # '000001000825',
                'REF_DOC_IT': '0001',
            })

        print("Added first record in paramsGOODSMOVEMENTS", paramsGOODSMOVEMENTS)

        for matr in sap_order:
            paramsGOODSMOVEMENTS.append({
                'MATERIAL': matr['MATNR'],  # '000000000002003411',
                'ENTRY_QNT': matr['ERFMG'],  # '1',
                'ENTRY_UOM': 'ST',
                'STGE_LOC': '7813',
                'BATCH': matr['CHARG'],  # '1000001747',
                'MOVE_TYPE': '261',
                'SPEC_STOCK': '',  # Индикатор особого запаса
                'WBS_ELEM': '',  # СПП-элемент
                'NO_MORE_GR': '',  # = 'X' если конечное подтверждение
                'RESERV_NO': matr['RSNUM'],  # '0000031904',  # Номер резерва
                'RES_ITEM': matr['RSPOS'],  # '0020',  # Номер позиции резерва
                'PLANT': 'SL00',
                'ORDERID': '00000' + SAP_ORDER,  # '000001000825',
                'WITHDRAWN': 'X',  # фиксированное значение
                'REF_DOC_IT': '0001',  # фиксированное значение
            })
        print("Added second and after records in paramsGOODSMOVEMENTS",
              paramsGOODSMOVEMENTS)
        # breakpoint()
        # for row in rowsinfoOrder:
        #     for matr in sap_order:  # infoMaterialOrder:
        #         paramsGOODSMOVEMENTS = [
        #             {
        #                 'MATERIAL': productsap,  # '000000000003100302',
        #                 'PLANT': 'SL00',
        #                 'MOVE_TYPE': '131',
        #                 'ENTRY_QNT': row['Qty'],  # '1',
        #                 'ENTRY_UOM': 'ST',
        #                 'ORDERID': '00000' + row['WO'],  # '000001000825',
        #                 'REF_DOC_IT': '0001',
        #             },
        #             {
        #                 'MATERIAL': matr['MATNR'],  # '000000000002003411',
        #                 'ENTRY_QNT': matr['ERFMG'],  # '1',
        #                 'ENTRY_UOM': 'ST',
        #                 'STGE_LOC': '7813',
        #                 'BATCH': matr['CHARG'],  # '1000001747',
        #                 'MOVE_TYPE': '261',
        #                 'SPEC_STOCK': '',  # Индикатор особого запаса
        #                 'WBS_ELEM': '',  # СПП-элемент
        #                 'NO_MORE_GR': '',  # = 'X' если конечное подтверждение
        #                 'RESERV_NO': matr['RSNUM'],  # '0000031904',  # Номер резерва
        #                 'RES_ITEM': matr['RSPOS'],  # '0020',  # Номер позиции резерва
        #                 'PLANT': 'SL00',
        #                 'ORDERID': '00000' + SAP_ORDER,  # '000001000825',
        #                 'WITHDRAWN': 'X',  # фиксированное значение
        #                 'REF_DOC_IT': '0001',  # фиксированное значение
        #             },
        #
        #         ]

        outputtedorder = connection.call('Z_IEXT_PRODORDCONF_CREATE_HDR', **{
            'UCODE': '21717',
            'PCODE': 'NEWPASSWORD1',
            'ATHDRLEVELS': paramsATHDRLEVELS,
            'GOODSMOVEMENTS': paramsGOODSMOVEMENTS,
        }
        )

        print(outputtedorder)

        parse_response(outputtedorder)
        logger.warning(parse_response(outputtedorder))
        print(parse_response(outputtedorder))

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


def parse_response(dict_value):
    for key, value in dict_value.items():
        if key == 'RETURN':
            if value[0].get('TYPE', '') == 'E':
                return 'Error: ' + str(value[0].get('MESSAGE', ''))
            elif value[0].get('TYPE', '') == 'I':
                return "Infomation: " + str(value[0].get('MESSAGE', ''))
            elif value[0].get('TYPE', '') == 'W':
                return "Warning: " + str(value[0].get('MESSAGE', ''))
            # elif value[0].get('NUMBER', '') == '469':
            # print("NOK")
            # elif value[0].get('NUMBER', '') == '100':
            # print("OK")
            else:
                return "Ответ не получен"


if __name__ == "__main__":
    main()
