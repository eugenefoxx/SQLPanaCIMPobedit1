import configparser

import pyrfc
import csv
from collections import defaultdict
from pyrfc import Connection

from pyrfc import ABAPApplicationError, ABAPRuntimeError, LogonError, CommunicationError, ExternalRuntimeError
import logging


def main():
    logger = logging.getLogger("order_info")
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
      #  result = connection.call('STFC_CONNECTION', REQUTEXT=u'Hello SAP!')
      #  print(result)
        logger.info("Connection to SAP RFC creating.")

        work_order = None
        # with open('/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data/work_order_name.csv', newline='') as csvfile:
        with open('/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test/1000862_work_order_name.csv', newline='') as csvfile:
            # with open('/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test/test1_work_order_name.csv', newline='') as csvfile:
            # with open('/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test/test2_work_order_name.csv', newline='') as csvfile:
            spamreader = csv.reader(csvfile, delimiter=',', quotechar='|')
            for row in spamreader:
                work_order = '' .join(row)

        orderSAP = '00000' + work_order  # 000001000825 000001000836
        # wo_component.csv wo_component_1000836.csv
        # 'wo_component_1000862.csv'
        # wo_component = '/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data/wo_component.csv'
        wo_component = '/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test/wo_component_1000862.csv'
        # wo_component = '/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test/test1_wo_component.csv'
        # wo_component = '/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test/test2_wo_component.csv'

        #infoOrder = '/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data/info_order.csv'
        infoOrder = '/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test/1000862_info_order.csv'
        # infoOrder = '/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test/test1_info_order.csv'
        # infoOrder = '/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test/test2_info_order.csv'

        order_info = connection.call('Z_IEXT_PRODORD_INFO', **{
            'AUFNR': str(orderSAP),  # 000001000825
            'UCODE': '21717',
            'PCODE': 'NEWPASSWORD1',
        }
        )
        ordp = order_info['RESITEMS']
        productsap = order_info['PRODUCT']
        print("productsap", productsap)
        print(order_info['RESITEMS'])
        print("_________inforesSAPorder___________")

        resSAPorder = [{'MATNR': sub['MATNR'], 'RSPOS': sub['RSPOS']}
                       for sub in ordp]

        # print(rrrr)
        print("type", type(resSAPorder))
        print("resSAPorder:", resSAPorder)
        for key in resSAPorder:
            if key['MATNR'] == '000000000001009223':
                print(key['RSPOS'])
            # print(key, '->', rrrr[key])
        print("_________info___________")
        print(order_info)

        print("reserv number:", get_reserv_num(order_info))
        reserv = get_reserv_num(order_info)
        # new order - 1000836
        # записываем в массив список компонентов в заказе
        resmtlrs = [sub['MATNR'] for sub in ordp]
        print("resmtlrs: " + str(resmtlrs))

        # valueProduct = '10'

        rowsinfoOrder = []
        with open(infoOrder, newline='') as file:
            csvreader = csv.DictReader(file, delimiter=',')
            for row in csvreader:
                rowsinfoOrder.append(row)
        print("rowsinfoOrder:", rowsinfoOrder)
        for c in rowsinfoOrder:
            print("Product:", productsap)
            print("Order", '00000' + c['WO'])
            print("Qty", c['Qty'])

        for pf in resmtlrs:
            if pf.__contains__('0000000000031'):
                print("31*", pf)
                # chg = connection.call('Z_IEXT_PRODORD_CHGRES', **{
                #    'UCODE': '21717',
                #    'PCODE': 'NEWPASSWORD1',
                #    'RESITEMS': {u'LINE_ID': 1, u'MATERIAL': pf, u'PLANT': 'SL00', u'STGE_LOC': '7813',
                #                 u'MOVE_TYPE': '261', u'ENTRY_QNT': row['SUM'], u'ENTRY_UOM': 'ST',
                #                 u'ORDERID': orderSAP, u'RESERV_NO': reserv, u'RES_ITEM': c['RSPOS']},
                # })
            if not pf.__contains__('0000000000031'):
                print("not have 31*")
                # break

        for i in resmtlrs:
            if i.__contains__('0000000000031'):
                print("i 31*: ", i)
            if not i.__contains__('0000000000031'):
                print("not i: ", i)

        for pcb in resmtlrs:
            if pcb.__contains__('200-'):
                print("pcb:", pcb)
            if not pcb.__contains__('200-'):
                print("not pcb:", pcb)
                break

        # убираем, если есть полуфабрикат из листа
        arrPRODORD_INFO_Component = []
        for pf in resmtlrs:
            if not pf.__contains__('0000000000031'):
                arrPRODORD_INFO_Component.append(pf)
        rowsPanaData = []
        with open(wo_component, newline='') as file:
            # csvreader = csv.reader(file, delimiter=',')
            csvreader = csv.DictReader(file, delimiter=',')
            # header = next(csvreader)
            for row in csvreader:
                rowsPanaData.append(row)
        print("rowsPanaData - ", rowsPanaData)
        for row in rowsPanaData:
            for c in resSAPorder:  # for c in resmtlrs:
                #    if c.__contains__('0000000000031'):
                #        print("31*", c)
                # проверка компонентов по Panacim на наличие в заказе
                if str('00000000000' + row['PART_NO']) == str(c['MATNR']):
                    #    if not c.__contains__('0000000000031'):
                    #    print('00000000000'+ row['PART_NO'])
                    print("This sap have in order", c['MATNR'], c['RSPOS'])

                    chg = connection.call('Z_IEXT_PRODORD_CHGRES', **{
                        'UCODE': '21717',
                        'PCODE': 'NEWPASSWORD1',
                        'RESITEMS': [
                            {
                                u'LINE_ID': '2',
                                u'MATERIAL': c['MATNR'],
                                u'PLANT': 'SL00',
                                u'STGE_LOC': '7813',
                                u'BATCH': row['Lot'],
                                u'MOVE_TYPE': '261',
                                u'ENTRY_QNT': row['SUM'],
                                u'ENTRY_UOM': 'ST',
                                u'ORDERID': orderSAP,
                                u'RESERV_NO': reserv,
                                u'RES_ITEM': c['RSPOS'],
                            }
                        ]
                    })
                    print(chg)
                    # logger.warning(parse_response(chg))
        #    elif str('00000000000' + row['PART_NO']) != str(c):
        #        print("not", c)

        # запись компонентов из Панасим в массив
        arrComponentFromPanaCIM = []
        with open(wo_component, newline='') as file:
            csvreader = csv.DictReader(file, delimiter=',')
            for row in csvreader:
                arrComponentFromPanaCIM.append('00000000000' + row['PART_NO'])
        print("arrComponentFromPanaCIM -", arrComponentFromPanaCIM)

        dict_of_lists_pana = defaultdict(list)
        for record in csv.DictReader(open(wo_component)):
            for key, val in record.items():
                dict_of_lists_pana[key].append(val)
        print("dict_of_lists_pana:", dict_of_lists_pana)

        # повторно читаем информацию по заказу для добавления на дробление сап и партии
        order_info = connection.call('Z_IEXT_PRODORD_INFO', **{
            'AUFNR': orderSAP,  # 000001000825
            'UCODE': '21717',
            'PCODE': 'NEWPASSWORD1',
        }
        )
        ordp = order_info['RESITEMS']

        # resSAPorder2 = [{'MATNR': sub['MATNR'], 'BDMNG': sub['BDMNG'], 'CHARG': sub['CHARG']} for sub in ordp]
        # resSAPorder2 = [{'PART_NO': sub['MATNR'], 'SUM': sub['BDMNG'], 'Lot': sub['CHARG']} for sub in ordp]

        resSAPorder2 = {}
        for i in ordp:
            resSAPorder2[i['MATNR']] = {"sum": i["BDMNG"], "lot": i["CHARG"]}
        # resSAPorder2 = [{'Lot': sub['CHARG']} for sub in ordp]
        print("resSAPorder2:", resSAPorder2)
        resarrComponentFromPanaCIM = []
        print("rowsPanaData: ", rowsPanaData)

        # добавлем в SAP заказ компонент, согласно PanaCIM, если его там нет
        for i in rowsPanaData:  # arrComponentFromPanaCIM
            # if '00000000000' + i['PART_NO'] and i['Lot'] not in resSAPorder2:
            # if ('00000000000' + i['PART_NO'] and i['SUM'] and i['Lot']) not in resSAPorder2:
            print(str('00000000000' + i["PART_NO"]))
            if str('00000000000' + i["PART_NO"]) in resSAPorder2.keys():
                #  добавляем позицию исходя из списка Панасим
                if i["Lot"] != resSAPorder2['00000000000' + i["PART_NO"]]["lot"]:
                    print(f"Lot {i['Lot']} {i['PART_NO']} fucked")
                    addcomp = connection.call('Z_IEXT_PRODORD_CHGRES', **{
                        'UCODE': '21717',
                        'PCODE': 'NEWPASSWORD1',
                        'RESITEMS': [
                            {
                                u'LINE_ID': '1',
                                u'MATERIAL': '00000000000' + i['PART_NO'],
                                u'PLANT': 'SL00',
                                u'STGE_LOC': '7813',
                                u'BATCH': i['Lot'],
                                u'MOVE_TYPE': '261',
                                u'ENTRY_QNT': i['SUM'],
                                u'ENTRY_UOM': 'ST',
                                u'ORDERID': orderSAP,
                                #    u'RESERV_NO': reserv,
                                #    u'RES_ITEM': c['RSPOS'],
                            }
                        ]
                    })
                    print(addcomp)
            else:
                # for c in resmtlrs:
                #  if '00000000000' + i['PART_NO'] != c['MATNR']:
                # print("component not have in sap_order:", i['PART_NO'], i['SUM'], i['Lot'])

                print(f"{i['PART_NO']} not found")
                # resarrComponentFromPanaCIM.append(i)

        print(resarrComponentFromPanaCIM)
        with open(wo_component, newline='') as file:
            csvreader = csv.DictReader(file, delimiter=',')
            for row in csvreader:
                for i in resarrComponentFromPanaCIM:
                    if i == '00000000000' + row['PART_NO']:
                        print("insert unknown sap in order", '00000000000' +
                              row['PART_NO'], '00000000000' + row['SUM'])

        print("resmtlrs", resmtlrs)

        # поиск компонентов, которые отсутствуют в cписке Panacim

        for i in resSAPorder:  # arrPRODORD_INFO_Component
            # print("remove:", str(i['MATNR']).removeprefix('00000000000'))
            # ch = str(i['MATNR']).removeprefix('00000000000')
            # print("ch:", ch)
            # arrComponentFromPanaCIM
            if i['MATNR'] not in arrComponentFromPanaCIM and not i['MATNR'].__contains__('0000000000031'):
                print("not in", i['MATNR'])
                zerocomp = connection.call('Z_IEXT_PRODORD_CHGRES', **{
                    'UCODE': '21717',
                    'PCODE': 'NEWPASSWORD1',
                    'RESITEMS': [
                        {
                            u'LINE_ID': '6',
                            # u'MATERIAL': i['MATNR'],
                            u'PLANT': 'SL00',
                            # u'STGE_LOC': '7813',
                            # u'MOVE_TYPE': '261',
                            # u'ENTRY_QNT': '0.0',
                            # u'ENTRY_UOM': 'ST',
                            u'ORDERID': orderSAP,
                            u'RESERV_NO': reserv,
                            u'RES_ITEM': i['RSPOS'],
                        }
                    ]
                })
                print(zerocomp)
                logger.info(zerocomp)

        order_info = connection.call('Z_IEXT_PRODORD_INFO', **{
            'AUFNR': orderSAP,  # 000001000825
            'UCODE': '21717',
            'PCODE': 'NEWPASSWORD1',
        }
        )
        #  Вставка партии в добавленные строки
        sap_order = order_info['RESITEMS']
        print("Состояние проверки измененного заказа", sap_order)
        # Не правильно тут делать ключ для парт-номера
        resSAPorder3 = {}
        for i in sap_order:
            # resSAPorder3[i['MATNR']] = {"sum": i["BDMNG"], "batch": i["CHARG"]}
            resSAPorder3[i['CHARG']] = {"sum": i["BDMNG"], "matrn": i["MATNR"]}
        # resSAPorder2 = [{'Lot': sub['CHARG']} for sub in ordp]
        print("resSAPorder3:", resSAPorder3)

        good_order = []
        bad_order = []
        rspos_black_list = []
        for i in sap_order:
            if not i['CHARG']:
                print(
                    f"Part number {i['MATNR']}, rspos {i['RSPOS']} does not have a lot")
                bad_order.append(i)
            else:
                good_order.append(i)
        print("bad_order", bad_order)
        print("rowsPanaData", rowsPanaData)
        for i in rowsPanaData:
            match = False
            for item in good_order:
                if item["MATNR"] == '00000000000' + i["PART_NO"] and item["CHARG"] == i["Lot"]:
                    match = True
                    break
                elif item["MATNR"] == '00000000000' + i["PART_NO"] and not item["CHARG"] == i["Lot"]:
                    match = False
            if not match:
                for item in bad_order:
                    if '00000000000' + i["PART_NO"] == item["MATNR"] and item['RSPOS'] not in rspos_black_list:
                        print({"Item of Bad_Order PART_NO": i['PART_NO'], "Lot": i['Lot'], "sum": i['SUM'],
                               "RSPOS": item['RSPOS']})
                        rspos_black_list.append(item['RSPOS'])

                        chg = connection.call('Z_IEXT_PRODORD_CHGRES', **{
                            'UCODE': '21717',
                            'PCODE': 'NEWPASSWORD1',
                            'RESITEMS': [
                                {
                                    u'LINE_ID': '2',
                                    u'MATERIAL': '00000000000' + i['PART_NO'],
                                    u'PLANT': 'SL00',
                                    u'STGE_LOC': '7813',
                                    u'BATCH': i['Lot'],
                                    u'MOVE_TYPE': '261',
                                    u'ENTRY_QNT': i['SUM'],
                                    u'ENTRY_UOM': 'ST',
                                    u'ORDERID': orderSAP,
                                    u'RESERV_NO': reserv,
                                    u'RES_ITEM': item['RSPOS'],
                                }
                            ]
                        })
                        print(chg)
                        break

                    elif item['RSPOS'] in rspos_black_list:
                        pass
                    else:
                        print(
                            f"Part number {i['PART_NO']} was not found in SAPOrder")
        # for item in ordp:  # rowsPanaData
        #     match = False
        #     for i in rowsPanaData:  # ordp
        #         if item['MATNR'] == str('00000000000' + i['PART_NO']) and item['CHARG'] == i['Lot']:
        #             match = True
        #             print('all good')
        #             break
        #         else:
        #             match = False
        #     if not match:
        #         print(f"cannot find a match for part_number {i['PART_NO']}, lot {i['Lot']}")
        #         chg = connection.call('Z_IEXT_PRODORD_CHGRES', **{
        #             'UCODE': '21717',
        #             'PCODE': 'NEWPASSWORD1',
        #             'RESITEMS': [
        #                 {
        #                     u'LINE_ID': '2',
        #                     u'MATERIAL': '00000000000' + i['PART_NO'],
        #                     u'PLANT': 'SL00',
        #                     u'STGE_LOC': '7813',
        #                     u'BATCH': i['Lot'],
        #                     u'MOVE_TYPE': '261',
        #                     u'ENTRY_QNT': i['SUM'],
        #                     u'ENTRY_UOM': 'ST',
        #                     u'ORDERID': orderSAP,
        #                     u'RESERV_NO': reserv,
        #                     u'RES_ITEM': item['RSPOS'],
        #                 }
        #             ]
        #         })
        #         print(chg)

        order_info = connection.call('Z_IEXT_PRODORD_INFO', **{
            'AUFNR': orderSAP,  # 000001000825
            'UCODE': '21717',
            'PCODE': 'NEWPASSWORD1',
        }
        )
        ordp = order_info['RESITEMS']
        print("Состояние проверки измененного заказа 2", ordp)

        #    for i in rowsPanaData:

        #        if str('00000000000' + i["PART_NO"]) in resSAPorder3.keys():
        #            if resSAPorder3[str('00000000000' + i["PART_NO"])]["batch"] == '':
        #                print("found not lot:", str('00000000000' + i["PART_NO"]), i['Lot'])
        #                print(f"{i['PART_NO']}, {i['lot']}")
        #                chg = connection.call('Z_IEXT_PRODORD_CHGRES', **{
        #                    'UCODE': '21717',
        #                    'PCODE': 'NEWPASSWORD1',
        #                    'RESITEMS': [
        #                        {
        #                            u'LINE_ID': '2',
        #                            u'MATERIAL': '00000000000' + i['PART_NO'],
        #                            u'PLANT': 'SL00',
        #                            u'STGE_LOC': '7813',
        #                            u'BATCH': i['Lot'],
        #                            u'MOVE_TYPE': '261',
        #                            u'ENTRY_QNT': i['SUM'],
        #                            u'ENTRY_UOM': 'ST',
        #                            u'ORDERID': orderSAP,
        #                            u'RESERV_NO': reserv,
        #                            u'RES_ITEM': resSAPorder3['RSPOS'],
        #                        }
        #                    ]
        #                })
        #                print(chg)

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


def get_reserv_num(dict_value):
    for key, value in dict_value.items():
        if key == 'RESITEMS':
            res = value[0].get('RSNUM', '')
            return res


def parse_response(dict_value):
    for value in dict_value['RETURN']:
        # if key == 'RETURN':
        if value.get('TYPE', '') == 'E':
            return 'Error: ' + str(value.get('MESSAGE', ''))
        elif value.get('TYPE', '') == 'I':
            return "Infomation: " + str(value.get('MESSAGE', ''))
        elif value.get('TYPE', '') == 'W':
            return "Warning: " + str(value.get('MESSAGE', ''))

        # elif value[0].get('NUMBER', '') == '469':
        # print("NOK")
        # elif value[0].get('NUMBER', '') == '100':
        # print("OK")
        else:
            return "Ответ не получен"


# def parse_response(dict_value):
#    for key, value in dict_value.items():
#        if key == 'RETURN':
#            if value[0].get('TYPE', '') == 'E':
#                return 'Error: ' + str(value[0].get('MESSAGE', ''))
#            elif value[0].get('TYPE', '') == 'I':
#                return "Infomation: " + str(value[0].get('MESSAGE', ''))
#            elif value[0].get('TYPE', '') == 'W':
#                return "Warning: " + str(value[0].get('MESSAGE', ''))
            # elif value[0].get('NUMBER', '') == '469':
            # print("NOK")
            # elif value[0].get('NUMBER', '') == '100':
            # print("OK")
#            else:
#                return "Ответ не получен"


if __name__ == "__main__":
    main()
