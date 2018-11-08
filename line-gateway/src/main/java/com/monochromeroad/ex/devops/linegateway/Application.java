/*
 * Copyright 2016 LINE Corporation
 *
 * LINE Corporation licenses this file to you under the Apache License,
 * version 2.0 (the "License"); you may not use this file except in compliance
 * with the License. You may obtain a copy of the License at:
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 */

package com.monochromeroad.ex.devops.linegateway;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

import java.util.Map;

import com.linecorp.bot.model.event.Event;
import com.linecorp.bot.model.event.MessageEvent;
import com.linecorp.bot.model.event.message.TextMessageContent;
import com.linecorp.bot.model.message.Message;
import com.linecorp.bot.model.message.TextMessage;
import com.linecorp.bot.spring.boot.annotation.EventMapping;
import com.linecorp.bot.spring.boot.annotation.LineMessageHandler;
import org.springframework.boot.web.client.RestTemplateBuilder;
import org.springframework.web.client.RestOperations;

@SpringBootApplication
@LineMessageHandler
public class Application {
    public static void main(String[] args) {
        SpringApplication.run(Application.class, args);
    }

    private final RestOperations restOperations;

    public Application(RestTemplateBuilder restTemplateBuilder) {
        this.restOperations = restTemplateBuilder.build();
    }

    @EventMapping
    public Message handleTextMessageEvent(MessageEvent<TextMessageContent> event) {
        System.out.println("event: " + event);
        final String originalMessageText = event.getMessage().getText().toLowerCase();

        String url = "http://devops-ex-greeter:8080/api/greeting";

        if (originalMessageText.contains("抽選")) {
            String number = originalMessageText.replaceAll("[^0-9]", "");
            url = "http://devops-ex-selector-knative.showcase-istio.svc.cluster.local/api/winners?number=" + number;
            Winners list = restOperations.getForObject(url, Winners.class);
            String message = "抽選結果です！！ \n" + String.join("\n", list.getWinners());
            return new TextMessage(message);
        }

        if (originalMessageText.toLowerCase().contains("knative")) {
            url = "http://devops-ex-greeter-knative.showcase-istio.svc.cluster.local/api/greeting";
        }

        Map<String, String> greeting = restOperations.getForObject(url, Map.class);
        return new TextMessage(greeting.get("content"));
    }

    @EventMapping
    public void handleDefaultMessageEvent(Event event) {
        System.out.println("event: " + event);
    }
}
