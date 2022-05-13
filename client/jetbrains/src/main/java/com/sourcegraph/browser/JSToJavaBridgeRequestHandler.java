package com.sourcegraph.browser;

import com.google.gson.JsonObject;
import com.intellij.openapi.project.Project;
import com.intellij.ui.jcef.JBCefJSQuery;
import com.sourcegraph.config.ConfigUtil;
import com.sourcegraph.config.ThemeUtil;

import javax.annotation.Nullable;

public class JSToJavaBridgeRequestHandler {
    private final Project project;

    public JSToJavaBridgeRequestHandler(Project project) {
        this.project = project;
    }

    public JBCefJSQuery.Response handle(JsonObject request) {
        String action = request.get("action").getAsString();
        // JsonObject arguments = request.getAsJsonObject("arguments");
        if (action.equals("getConfig")) {
            JsonObject configAsJson = new JsonObject();
            configAsJson.addProperty("instanceURL", ConfigUtil.getSourcegraphUrl(this.project));
            return createResponse(configAsJson);
        } else if (action.equals("getTheme")) {
            JsonObject currentThemeAsJson = ThemeUtil.getCurrentThemeAsJson();
            return createResponse(currentThemeAsJson);
        } else {
            return createResponse(2, "Unknown action: " + action, null);
        }
    }

    public JBCefJSQuery.Response handleInvalidRequest() {
        return createResponse(1, "Invalid JSON passed to bridge.", null);
    }

    private JBCefJSQuery.Response createResponse(@Nullable JsonObject result) {
        return new JBCefJSQuery.Response(result != null ? result.toString() : null);
    }

    private JBCefJSQuery.Response createResponse(int errorCode, @Nullable String errorMessage, @Nullable JsonObject data) {
        return new JBCefJSQuery.Response(data != null ? data.toString() : null, errorCode, errorMessage);
    }
}
