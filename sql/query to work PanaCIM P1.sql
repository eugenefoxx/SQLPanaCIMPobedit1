
/* поиск последнего закрывшегося job_id */
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
  
 /* результат - job_id and setup_id */
  
 /* узнаем кол-во м / з произведенных  */
  /****** Script for SelectTopNRows command from SSMS  ******/
SELECT COUNT(DISTINCT PANEL_ID) 
  FROM [PanaCIM].[dbo].[panels] where JOB_ID = '5134'
  
 /* результат1 - 250 */
  
/*  получаем номер NC_VERSION (одной заготовки) - */
  /****** Script for SelectTopNRows command from SSMS  ******/
SELECT top 1 [NC_VERSION]
FROM [PanaCIM].[dbo].[panels] where JOB_ID = '5133' order by NC_VERSION desc

/* результат 298118 */

/* получем кол-во плат в заготовке */
SELECT MAX(PATTERN_NUMBER)
  FROM [PanaCIM].[dbo].[nc_placement_detail] where NC_VERSION = '298118'
/* результат2 - 8 */
*/
2000 плат = результат1 - 250 * результат2 - 8
*/
/* по компонентам -  
узнать уникальные номера панел ид  */
/****** Script for SelectTopNRows command from SSMS  ******/
SELECT DISTINCT PANEL_ID
FROM [PanaCIM].[dbo].[panels] where JOB_ID = '5134' /*order by NC_VERSION desc*/

/* получить список NC_VERSION */
/****** Script for SelectTopNRows command from SSMS  ******/
SELECT DISTINCT NC_VERSION
  FROM [PanaCIM].[dbo].[panels] where JOB_ID = '5134'
  
-- Получаем список REPORT_ID
  /****** Script for SelectTopNRows command from SSMS  ******/
SELECT TOP 1000 [REPORT_ID]
  FROM [PanaCIM].[dbo].[PRODUCTION_REPORTS_NM_NPM_KX_VIEW]
where NC_VERSION = '297832'
    and JOB_ID = '5134'
    AND PLACE_COUNT != '0'
  
/*  Информация по потреблению компонентов по REPORT_ID  */
  /****** Script for SelectTopNRows command from SSMS  ******/
SELECT [REPORT_ID]
      ,[LOT]
      ,[STAGE]
      ,[TPROD]
      ,[HEAD]
      ,[FADD]
      ,[FSADD]
      ,[FBLKCODE]
      ,[FBLKSERIAL]
      ,[REELID]
      ,[TCNT]
      ,[TMISS]
      ,[RMISS]
      ,[HMISS]
      ,[FMISS]
      ,[MMISS]
      ,[BOARD]
      ,[PARTSNAME]
      ,[PLACE_COUNT]
      ,[TRSMISS]
      ,[TIMESTAMP]
      ,[REEL_ID]
      ,[FEEDER_ID]
      ,[PPIERR_COUNT]
  FROM [PanaCIM].[dbo].[Z_Cass_NPM] where REPORT_ID = '8211776' order by REPORT_ID desc
  
  /* Вариант 2 - смотреть потребление через VIEW */
  /****** Script for SelectTopNRows command from SSMS  ******/
SELECT TOP 10000 SUM(PLACE_COUNT) AS SUM/*, [REPORT_ID]
      ,[PRODUCT_ID]
      ,[EQUIPMENT_ID]
      ,[START_TIME]
      ,[END_TIME]
      ,[SETUP_ID]
      ,[NC_VERSION]
      ,[LANE_NO]
      ,[JOB_ID]
      ,[CASSETTE]
      ,[COMP_ID]
      ,[REEL_ID]
      ,[FEEDER_BARCODE]
      ,[PICKUP_COUNT]
      ,[PLACE_COUNT]
      ,[PICKUP_MISS_COUNT]
      ,[PICKUP_ERROR_COUNT]
      ,[RECOG_ERROR_COUNT]
      ,[HEIGHT_MISS_COUNT]
      ,[SHAPE_ERROR_COUNT]
      ,[DROP_MISS_COUNT]
      ,[MOUNT_MISS_COUNT]
      ,[TABLE_NO]
      ,[TABLE_SLOT]
      ,[TABLE_SUBSLOT]
      ,[PU_NUMBER]
      ,[VENDOR_NO]
      ,[LOT_NO]
      ,[FEEDER_ID]
      ,[TRANSFER_UNIT_PART_DROP_ERROR_COUNT]
      ,[SHAPE_CODE]*/
  FROM [PanaCIM].[dbo].[Z_CASS_VIEW] where JOB_ID = '5109' AND REEL_ID = '213953' group by REEL_ID
  /* Вариант 2-b - ЕО и кол-во */
  /****** Script for SelectTopNRows command from SSMS  ******/
SELECT TOP 10000 [REEL_ID],SUM(PLACE_COUNT) AS SUM
  FROM [PanaCIM].[dbo].[Z_CASS_VIEW] where JOB_ID = '5109' group by REEL_ID

/*  Вариант 2-d - ЕО и кол-во  */
/****** Script for SelectTopNRows command from SSMS  ******/
SELECT TOP 10000 [PanaCIM].[dbo].[Z_CASS_VIEW].[REEL_ID],SUM([PanaCIM].[dbo].[Z_CASS_VIEW].PLACE_COUNT) AS SUM, [PanaCIM].[dbo].[REEL_DATA_VIEW].reel_barcode
  FROM [PanaCIM].[dbo].[Z_CASS_VIEW]LEFT JOIN [PanaCIM].[dbo].[REEL_DATA_VIEW] 
  ON [PanaCIM].[dbo].[REEL_DATA_VIEW].[reel_id] = [PanaCIM].[dbo].[Z_CASS_VIEW].[REEL_ID] 
  where [PanaCIM].[dbo].[Z_CASS_VIEW].JOB_ID = '5109' 
  group by [PanaCIM].[dbo].[REEL_DATA_VIEW].reel_barcode, [PanaCIM].[dbo].[Z_CASS_VIEW].REEL_ID

/****** Script for SelectTopNRows command from SSMS  ******/
SELECT TOP 10000 [PanaCIM].[dbo].[Z_CASS_VIEW].[REEL_ID],SUM([PanaCIM].[dbo].[Z_CASS_VIEW].PLACE_COUNT) AS SUM, [PanaCIM].[dbo].[REEL_DATA_VIEW].reel_barcode
  FROM [PanaCIM].[dbo].[Z_CASS_VIEW]LEFT JOIN [PanaCIM].[dbo].[REEL_DATA_VIEW]
  ON [PanaCIM].[dbo].[REEL_DATA_VIEW].[reel_id] = [PanaCIM].[dbo].[Z_CASS_VIEW].[REEL_ID]
  where [PanaCIM].[dbo].[Z_CASS_VIEW].JOB_ID = '5109' 
  group by [PanaCIM].[dbo].[REEL_DATA_VIEW].reel_barcode, [PanaCIM].[dbo].[Z_CASS_VIEW].REEL_ID
  order by [PanaCIM].dbo.REEL_DATA_VIEW.reel_barcode desc
