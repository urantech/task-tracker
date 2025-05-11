package com.urantech.reportservice;

import com.urantech.reportservice.service.ReportService;
import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequiredArgsConstructor
@RequestMapping("/api/admin/reports")
public class ReportController {

    private final ReportService reportService;

    @GetMapping
    public void generateAndSendReports() {
        reportService.generateAndSendReports();
    }
}
