package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/cli"
)

const (
	version = "0.0.0"
)

func loadCommands() map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"ability list":  AbilityListCommand,
		"ability test":  AbilityTestCommand,
		"addon list":    AddonListCommand,
		"addon install": AddonInstallCommand,
		"addon show":    AddonShowCommand,
		"addon delete":  AddonDeleteCommand,
		"addon update":  AddonUpdateCommand,

		"escalation-policy list":   EscalationPolicyListCommand,
		"escalation-policy create": EscalationPolicyCreateCommand,
		"escalation-policy delete": EscalationPolicyDeleteCommand,
		"escalation-policy show":   EscalationPolicyShowCommand,
		"escalation-policy update": EscalationPolicyUpdateCommand,

		"event-orchestration list":   EventOrchestrationListCommand,
		"event-orchestration create": EventOrchestrationCreateCommand,
		"event-orchestration delete": EventOrchestrationDeleteCommand,
		"event-orchestration show":   EventOrchestrationShowCommand,
		"event-orchestration update": EventOrchestrationUpdateCommand,

		"event-v2 manage": EventV2ManageCommand,

		"incident list":        IncidentListCommand,
		"incident manage":      IncidentManageCommand,
		"incident show":        IncidentShowCommand,
		"incident note list":   IncidentNoteListCommand,
		"incident note create": IncidentNoteCreateCommand,
		"incident snooze":      IncidentSnoozeCommand,

		"log-entry list": LogEntryListCommand,
		"log-entry show": LogEntryShowCommand,

		"maintenance-window list":   MaintenanceWindowListCommand,
		"maintenance-window create": MaintenanceWindowCreateCommand,
		"maintenance-window delete": MaintenanceWindowDeleteCommand,
		"maintenance-window show":   MaintenanceWindowShowCommand,
		"maintenance-window update": MaintenanceWindowUpdateCommand,

		"notification list": NotificationListCommand,

		"oncall list": OncallListCommand,

		"schedule list":    ScheduleListCommand,
		"schedule create":  ScheduleCreateCommand,
		"schedule preview": SchedulePreviewCommand,
		"schedule delete":  ScheduleDeleteCommand,
		"schedule show":    ScheduleShowCommand,
		"schedule update":  ScheduleUpdateCommand,

		"schedule override list":   ScheduleOverrideListCommand,
		"schedule override create": ScheduleOverrideCreateCommand,
		"schedule override delete": ScheduleOverrideDeleteCommand,

		"schedule oncall list": ScheduleOncallListCommand,

		"service list":               ServiceListCommand,
		"service create":             ServiceCreateCommand,
		"service delete":             ServiceDeleteCommand,
		"service show":               ServiceShowCommand,
		"service update":             ServiceUpdateCommand,
		"service integration create": ServiceIntegrationCreateCommand,
		"service integration show":   ServiceIntegrationShowCommand,
		"service integration update": ServiceIntegrationUpdateCommand,
		"service rule create":        ServiceRuleCreateCommand,
		"service rule delete":        ServiceRuleDeleteCommand,
		"service rule list":          ServiceRuleListCommand,
		"service rule show":          ServiceRuleShowCommand,
		"service rule update":        ServiceRuleUpdateCommand,

		"team list":                     TeamListCommand,
		"team create":                   TeamShowCommand,
		"team delete":                   TeamDeleteCommand,
		"team show":                     TeamShowCommand,
		"team update":                   TeamUpdateCommand,
		"team remove escalation-policy": TeamRemoveEscalationPolicyCommand,
		"team add escalation-policy":    TeamAddEscalationPolicyCommand,
		"team add user":                 TeamAddUserCommand,

		"user list":                     UserListCommand,
		"user create":                   UserCreateCommand,
		"user delete":                   UserDeleteCommand,
		"user show":                     UserShowCommand,
		"user update":                   UserUpdateCommand,
		"user contact-method list":      UserContactMethodListCommand,
		"user contact-method create":    UserContactMethodCreateCommand,
		"user contact-method delete":    UserContactMethodDeleteCommand,
		"user contact-method show":      UserContactMethodShowCommand,
		"user contact-method update":    UserContactMethodUpdateCommand,
		"user notification-rule list":   UserNotificationRuleListCommand,
		"user notification-rule create": UserNotificationRuleCreateCommand,
		"user notification-rule delete": UserNotificationRuleDeleteCommand,
		"user notification-rule show":   UserNotificationRuleShowCommand,
		"user notification-rule update": UserNotificationRuleUpdateCommand,

		"vendor list":             VendorListCommand,
		"vendor show":             VendorShowCommand,
		"analytics incident show": AnalyticsGetAggregatedIncidentDataCommand,
		"analytics service show":  AnalyticsGetAggregatedServiceDataCommand,
		"analytics team show":     AnalyticsGetAggregatedTeamDataCommand,
	}
}

func main() {
	os.Exit(invokeCLI(os.Args[1:]))
}

func invokeCLI(args []string) int {
	for _, arg := range args {
		if arg == "-v" || arg == "--version" {
			fmt.Println(version)
			return 0
		}
	}

	cli := &cli.CLI{
		Args:     args,
		Commands: loadCommands(),
		HelpFunc: cli.BasicHelpFunc("pd"),
	}

	exitCode, err := cli.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return 1
	}

	return exitCode
}
