package com.urantech.reportservice.service;

import lombok.RequiredArgsConstructor;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;

@Component
@RequiredArgsConstructor
public class ReportSender {

    private final ReportService reportService;

    @Scheduled(cron = "${spring.scheduler.cron}")
    public void generateDailyReport() {
        reportService.generateAndSendReports();
    }
}
