/* Поиск колонки в базе */
SELECT      c.name  AS 'ColumnName'
            ,t.name AS 'TableName'
FROM        sys.columns c
JOIN        sys.tables  t   ON c.object_id = t.object_id
WHERE       c.name LIKE '%idnum%' /*'%pattern%'   */
ORDER BY    TableName
            ,ColumnName;

/* По ворк-ордеру определяем job_id */
/****** Script for SelectTopNRows command from SSMS  ******/
SELECT TOP 1000 [WORK_ORDER_ID]
      ,[WORK_ORDER_NAME]
      ,[LOT_SIZE]
      ,[JOB_ID]
      ,[MASTER_WORK_ORDER_ID]
      ,[COMMENTS]
  FROM [PanaCIM].[dbo].[work_orders]where WORK_ORDER_NAME = 'No5897_TM11_pri'
  /*Результат*.
  WORK_ORDER_ID,WORK_ORDER_NAME,LOT_SIZE,JOB_ID,MASTER_WORK_ORDER_ID,COMMENTS
6621,No5897_TM11_pri,3000,5131,-1,NULL

/* Узнаем job_id и по нему проводим расчет произведенных м/з*/
/****** Script for SelectTopNRows command from SSMS  ******/
SELECT COUNT(DISTINCT PANEL_ID) /*TOP 1000 [PANEL_ID]
      ,[EQUIPMENT_ID]
      ,[NC_VERSION]
      ,[START_TIME]
      ,[END_TIME]
      ,[PANEL_EQUIPMENT_ID]
      ,[PANEL_SOURCE]
      ,[PANEL_TRACE]
      ,[STAGE_NO]
      ,[LANE_NO]
      ,[JOB_ID]
      ,[SETUP_ID]
      ,[TRX_PRODUCT_ID]
      ,[HUMIDITY]
      ,[TEMPERATURE]
      ,[SETUP_AUDIT_HISTORY_ID]
      ,[AUDIT_STATUS]*/
  FROM [PanaCIM].[dbo].[panels] where JOB_ID = '5131'
  Результат - 750 
  
  /* Узнать кол-во плат в заготовке */
  /****** Script for SelectTopNRows command from SSMS  ******/
SELECT MAX(PATTERN_NUMBER)/*COUNT(PATTERN_NUMBER)*/ /*TOP 1000 [NC_PLACEMENT_ID]
      ,[IDNUM]
      ,[REF_DESIGNATOR]
      ,[PART_NAME]
      ,[PATTERN_NUMBER]
      ,[NC_VERSION]
      ,[PATTERN_IDNUM]
      ,[PU_DISPLAY]*/
  FROM [PanaCIM].[dbo].[nc_placement_detail] where NC_VERSION = '298113' /*AND PATTERN_NUMBER = '1'*/
  
  /* Просмотр закрытия по времени ворк-ордеров, если CLOSING_TYPE 0 то закрыт, если NULL - то в производстве    */
  /****** Script for SelectTopNRows command from SSMS  ******/
SELECT TOP 1000 [JOB_ID]
      ,[EQUIPMENT_ID]
      ,[SETUP_ID]
      ,[START_TIME]
      ,[END_TIME]
      ,[CLOSING_TYPE]
      ,[START_OPERATOR_ID]
      ,[END_OPERATOR_ID]
      ,[TFR_REASON]
      ,[LANE_NO]
  FROM [PanaCIM].[dbo].[job_history] order by END_TIME desc
  
  /****** Script for SelectTopNRows command from SSMS  ******/
SELECT TOP 1 [JOB_ID]
      ,[EQUIPMENT_ID]
      ,[SETUP_ID]
      ,[START_TIME]
      ,[END_TIME]
      ,[CLOSING_TYPE]
      ,[START_OPERATOR_ID]
      ,[END_OPERATOR_ID]
      ,[TFR_REASON]
      ,[LANE_NO]
  FROM [PanaCIM].[dbo].[job_history] where CLOSING_TYPE = '0' order by END_TIME desc
