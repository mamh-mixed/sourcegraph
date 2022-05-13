package com.sourcegraph.browser;

import com.google.gson.JsonObject;
import com.intellij.openapi.project.Project;
import com.intellij.ui.jcef.JBCefJSQuery;
import com.sourcegraph.config.ConfigUtil;
import com.sourcegraph.config.ThemeUtil;
import com.sourcegraph.find.Search;
import org.jetbrains.annotations.NotNull;

import javax.annotation.Nullable;

public class JSToJavaBridgeRequestHandler {
    private final Project project;

    public JSToJavaBridgeRequestHandler(@NotNull Project project) {
        this.project = project;
    }

    public JBCefJSQuery.Response handle(@NotNull JsonObject request) {
        String action = request.get("action").getAsString();
        JsonObject arguments = request.getAsJsonObject("arguments");
        switch (action) {
            case "getConfig": {
                try {
                    JsonObject configAsJson = new JsonObject();
                    configAsJson.addProperty("instanceURL", ConfigUtil.getSourcegraphUrl(this.project));
                    return createSuccessResponse(configAsJson);
                } catch (Exception e) {
                    return createErrorResponse(3, e.getClass().getName() + ": " + e.getMessage());
                }
            }
            case "getTheme":
                try {
                    JsonObject currentThemeAsJson = ThemeUtil.getCurrentThemeAsJson();
                    return createSuccessResponse(currentThemeAsJson);
                } catch (Exception e) {
                    return createErrorResponse(4, e.getClass().getName() + ": " + e.getMessage());
                }
            case "saveLastSearch":
                try {
                    String query = arguments.get("query").getAsString();
                    boolean caseSensitive = arguments.get("caseSensitive").getAsBoolean();
                    String patternType = arguments.get("patternType").getAsString();
                    String selectedSearchContextSpec = arguments.get("selectedSearchContextSpec").getAsString();
                    ConfigUtil.setLastSearch(project, new Search(
                        query,
                        caseSensitive,
                        patternType,
                        selectedSearchContextSpec
                    ));
                    return createSuccessResponse(new JsonObject());
                } catch (Exception e) {
                    return createErrorResponse(5, e.getClass().getName() + ": " + e.getMessage());
                }
            case "loadLastSearch": {
                try {
                    JsonObject configAsJson = new JsonObject();
                    Search lastSearch = ConfigUtil.getLastSearch(this.project);
                    configAsJson.addProperty("query", lastSearch.getQuery());
                    configAsJson.addProperty("caseSensitive", lastSearch.isCaseSensitive());
                    configAsJson.addProperty("patternType", lastSearch.getPatternType());
                    configAsJson.addProperty("selectedSearchContextSpec", lastSearch.getSelectedSearchContextSpec());
                    return createSuccessResponse(configAsJson);
                } catch (Exception e) {
                    return createErrorResponse(6, e.getClass().getName() + ": " + e.getMessage());
                }
            }
            default:
                return createErrorResponse(2, "Unknown action: " + action);
        }
    }

    public JBCefJSQuery.Response handleInvalidRequest() {
        return createErrorResponse(1, "Invalid JSON passed to bridge.");
    }

    @NotNull
    private JBCefJSQuery.Response createSuccessResponse(@Nullable JsonObject result) {
        return new JBCefJSQuery.Response(result != null ? result.toString() : null);
    }

    @NotNull
    private JBCefJSQuery.Response createErrorResponse(int errorCode, @Nullable String errorMessage) {
        return new JBCefJSQuery.Response(null, errorCode, errorMessage);
    }
}
