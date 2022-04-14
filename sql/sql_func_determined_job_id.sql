/* ?? ????-?????? ?????????? job_id */
/****** Script for SelectTopNRows command from SSMS  ******/
SELECT TOP 1000 [WORK_ORDER_ID]
      ,[WORK_ORDER_NAME]
      ,[LOT_SIZE]
      ,[JOB_ID]
      ,[MASTER_WORK_ORDER_ID]
      ,[COMMENTS]
FROM [PanaCIM].[dbo].[work_orders]
where JOB_ID = '5236'

/* ?? ????-?????? ?????????? job_id */
/****** Script for SelectTopNRows command from SSMS  ******/
SELECT TOP 1000 [WORK_ORDER_ID]
      ,[WORK_ORDER_NAME]
      ,[LOT_SIZE]
      ,[JOB_ID]
      ,[MASTER_WORK_ORDER_ID]
      ,[COMMENTS]
FROM [PanaCIM].[dbo].[work_orders] where WORK_ORDER_NAME = 'No6409_M17_PRI'

    /* ?????? ???-?? ? / ? ?????????????  */
--SELECT * FROM dbo.GetLastJobId()
    IF OBJECT_ID('dbo.SUMPattern') IS NOT NULL DROP FUNCTION  dbo.SUMPattern;
GO

CREATE FUNCTION dbo.SUMPattern()
    RETURNS TABLE
    AS
      --BEGIN
      RETURN
SELECT COUNT(DISTINCT PANEL_ID)AS sumPattern
FROM [PanaCIM].[dbo].[panels] where JOB_ID = '5235'/*(
            SELECT * FROM --dbo.GetLastJobId()) --'5134'
            --RETURN COUNT(DISTINCT PANEL_ID)*/
    -- END;
    --  END;
GO

/*  ???????? ????? NC_VERSION (????? ?????????) - */
    /****** Script for SelectTopNRows command from SSMS  ******/
    IF OBJECT_ID('dbo.GetNcVersion') IS NOT NULL DROP FUNCTION dbo.GetNcVersion;
GO

CREATE FUNCTION dbo.GetNcVersion()
    RETURNS TABLE
    AS
      RETURN
SELECT top 1 [NC_VERSION]
FROM [PanaCIM].[dbo].[panels] where JOB_ID = '5235' order by NC_VERSION desc
GO

/* ??????? ???-?? ???? ? ????????? */
    IF OBJECT_ID('dbo.CountPCBInPattern') IS NOT NULL DROP FUNCTION dbo.CountPCBInPattern;
GO

CREATE FUNCTION dbo.CountPCBInPattern()
    RETURNS TABLE
    AS
      RETURN
SELECT MAX(PATTERN_NUMBER) AS sumPCBInPattern
FROM [PanaCIM].[dbo].[nc_placement_detail] where NC_VERSION = (SELECT * FROM dbo.GetNcVersion()) --'298118'
GO

/* ??????? ????? ????????????? ???? */
    IF OBJECT_ID('dbo.SumProductionPCB') IS NOT NULL DROP FUNCTION dbo.SumProductionPCB;
GO

CREATE FUNCTION dbo.SumProductionPCB()
    RETURNS INT
BEGIN
RETURN (SELECT * FROM dbo.SUMPattern()) * (SELECT * FROM dbo.CountPCBInPattern())
END;
GO

/*  Вариант 2-d - ЕО и кол-во  */
/****** Script for SelectTopNRows command from SSMS  ******/
--IF OBJECT_ID('dbo.SumInstallComponent') IS NOT NULL DROP FUNCTION dbo.SumInstallComponent;
IF OBJECT_ID('dbo.SumPLACE_COUNT_ALL_REEL_ID') IS NOT NULL DROP FUNCTION dbo.SumPLACE_COUNT_ALL_REEL_ID;
GO

CREATE FUNCTION dbo.SumPLACE_COUNT_ALL_REEL_ID
(
    @reelid AS INT
--	@lastjobid AS INT
)
    RETURNS TABLE
    AS
RETURN
SELECT
--[PanaCIM].[dbo].[Z_CASS_VIEW].[REEL_ID],
SUM([PanaCIM].[dbo].[Z_CASS_VIEW].PLACE_COUNT) AS PLACE_COUNT
--SUM([PanaCIM].[dbo].[Z_CASS_VIEW].PICKUP_COUNT) AS PICKUP_COUNT
--[PanaCIM].[dbo].[REEL_DATA_VIEW].reel_barcode
FROM [PanaCIM].[dbo].[Z_CASS_VIEW] LEFT JOIN [PanaCIM].[dbo].[REEL_DATA_VIEW]
ON [PanaCIM].[dbo].[REEL_DATA_VIEW].[reel_id] = [PanaCIM].[dbo].[Z_CASS_VIEW].[REEL_ID]
    -- where [PanaCIM].[dbo].[Z_CASS_VIEW].JOB_ID = (SELECT * FROM dbo.GetLastJobId())
Where [PanaCIM].[dbo].[REEL_DATA_VIEW].[reel_id] = @reelid
group by [PanaCIM].[dbo].[REEL_DATA_VIEW].reel_barcode, [PanaCIM].[dbo].[Z_CASS_VIEW].REEL_ID
GO

    -- IF OBJECT_ID('dbo.SumInstallComponent') IS NOT NULL DROP FUNCTION dbo.SumInstallComponent;
    IF OBJECT_ID('dbo.SumPICKUP_COUNT_ALL_REEL_ID') IS NOT NULL DROP FUNCTION dbo.SumPICKUP_COUNT_ALL_REEL_ID;
GO

CREATE FUNCTION dbo.SumPICKUP_COUNT_ALL_REEL_ID
(
    @reelid AS INT
--	@lastjobid AS INT
)
    RETURNS TABLE
    AS
RETURN
SELECT
--[PanaCIM].[dbo].[Z_CASS_VIEW].[REEL_ID],
--SUM([PanaCIM].[dbo].[Z_CASS_VIEW].PLACE_COUNT) AS PLACE_COUNT
SUM([PanaCIM].[dbo].[Z_CASS_VIEW].PICKUP_COUNT) AS PICKUP_COUNT_ALL
--[PanaCIM].[dbo].[REEL_DATA_VIEW].reel_barcode
FROM [PanaCIM].[dbo].[Z_CASS_VIEW] LEFT JOIN [PanaCIM].[dbo].[REEL_DATA_VIEW]
ON [PanaCIM].[dbo].[REEL_DATA_VIEW].[reel_id] = [PanaCIM].[dbo].[Z_CASS_VIEW].[REEL_ID]
    -- where [PanaCIM].[dbo].[Z_CASS_VIEW].JOB_ID = (SELECT * FROM dbo.GetLastJobId())
Where [PanaCIM].[dbo].[REEL_DATA_VIEW].[reel_id] = @reelid
group by [PanaCIM].[dbo].[REEL_DATA_VIEW].reel_barcode, [PanaCIM].[dbo].[Z_CASS_VIEW].REEL_ID
GO

/****** Script for SelectTopNRows command from SSMS  ******/
    IF OBJECT_ID('dbo.InfoInstallLastJobId_View', 'V') IS NOT NULL DROP VIEW dbo.InfoInstallLastJobId_View
GO

CREATE VIEW dbo.InfoInstallLastJobId_View
AS
SELECT TOP 10000 [PanaCIM].[dbo].[Z_CASS_VIEW].[REEL_ID],
[PanaCIM].[dbo].[reel_data].PART_NO,
SUM([PanaCIM].[dbo].[Z_CASS_VIEW].PLACE_COUNT) AS PLACE_COUNT,
SUM([PanaCIM].[dbo].[Z_CASS_VIEW].PICKUP_COUNT) AS PICKUP_COUNT,
[PanaCIM].[dbo].[REEL_DATA_VIEW].reel_barcode,
[PanaCIM].[dbo].[reel_data].CURRENT_QUANTITY,
[PanaCIM].[dbo].[reel_data].QUANTITY AS INITIAL_QUANTITY
--(SELECT * FROM dbo.SumInstallComponent([PanaCIM].[dbo].[Z_CASS_VIEW].[REEL_ID]))
  FROM [PanaCIM].[dbo].[Z_CASS_VIEW]
 -- LEFT JOIN ( SELECT * FROM dbo.SumInstallComponent([PanaCIM].[dbo].[Z_CASS_VIEW].[REEL_ID]))
  LEFT JOIN [PanaCIM].[dbo].[REEL_DATA_VIEW]
  ON [PanaCIM].[dbo].[REEL_DATA_VIEW].[reel_id] = [PanaCIM].[dbo].[Z_CASS_VIEW].[REEL_ID]
  LEFT JOIN [PanaCIM].[dbo].[reel_data]
  ON [PanaCIM].[dbo].[REEL_DATA_VIEW].[reel_id] = [PanaCIM].[dbo].[reel_data].REEL_ID
  --where [PanaCIM].[dbo].[Z_CASS_VIEW].JOB_ID = (SELECT * FROM dbo.GetLastJobId()) AND [PanaCIM].[dbo].[Z_CASS_VIEW].[REEL_ID] IS NOT NULL
  where [PanaCIM].[dbo].[Z_CASS_VIEW].JOB_ID = '5236' AND [PanaCIM].[dbo].[Z_CASS_VIEW].[REEL_ID] IS NOT NULL
  group by [PanaCIM].[dbo].[REEL_DATA_VIEW].reel_barcode, [PanaCIM].[dbo].[Z_CASS_VIEW].REEL_ID, [PanaCIM].[dbo].[reel_data].CURRENT_QUANTITY, [PanaCIM].[dbo].[reel_data].QUANTITY,
  [PanaCIM].[dbo].[reel_data].PART_NO
--order by [PanaCIM].dbo.REEL_DATA_VIEW.reel_barcode desc
GO
--SELECT * FROM dbo.SumInstallComponent('218847');
--SELECT * FROM dbo.InfoInstallLastJobId_View;
SELECT TOP 1000 [REEL_ID]
	,[PART_NO]
      ,[PLACE_COUNT]
      ,[PICKUP_COUNT]
      ,[reel_barcode]
      ,[CURRENT_QUANTITY]
      ,[INITIAL_QUANTITY]
      , (SELECT * FROM dbo.SumPLACE_COUNT_ALL_REEL_ID([PanaCIM].[dbo].[InfoInstallLastJobId_View].[REEL_ID])) AS PLACE_COUNT_ALL
      , (SELECT * FROM dbo.SumPICKUP_COUNT_ALL_REEL_ID([PanaCIM].[dbo].[InfoInstallLastJobId_View].[REEL_ID])) AS PICKUP_COUNT_ALL
FROM [PanaCIM].[dbo].[InfoInstallLastJobId_View]
order by PART_NO;


--SELECT * FROM dbo.GetLastJobId();
SELECT * FROM dbo.SUMPattern();
SELECT * FROM dbo.GetNcVersion();
SELECT * FROM dbo.CountPCBInPattern();
PRINT dbo.SumProductionPCB();