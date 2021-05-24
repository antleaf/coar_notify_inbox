package main

//func testExpandJsonFromNotificationRecord(notification Notification) error {
//	var err error
//	jsonLdProcessor := ld.NewJsonLdProcessor()
//	var jsonLdProcessorOptions = ld.NewJsonLdOptions("")
//	_, err = jsonLdProcessor.Expand(notification.PayloadStruct, jsonLdProcessorOptions)
//	if err != nil {
//		zapLogger.Error(err.Error())
//		return err
//	}
//	return err
//}
//
//func testExpandJsonFromDbId(id uuid.UUID) error {
//	var err error
//	notification := LoadNotificationFromDbById(id)
//	err = testExpandJsonFromNotificationRecord(notification)
//	return err
//}

//func testxpandJsonFromFile(filePath string) error {
//	var err error
//	err = ioutil.ReadFile(filePath)
//	if err != nil {
//		zapLogger.Error(err.Error())
//		return err
//	}
//	notification := NewNotification("", time.Now())
//	err = testExpandJsonFromNotificationRecord(*notification)
//	return err
//}
