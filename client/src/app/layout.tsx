// app/layout.tsx

import { ToastProvider } from "@/contexts";
import "../styles/globals.scss";
import type { Metadata } from "next";
import { Manrope, Edu_NSW_ACT_Foundation } from "next/font/google";

const manrope = Manrope({
    subsets: ["latin"],
    weight: ["400", "500", "600", "700"],
    variable: "--font-manrope",
    display: "swap",
});

// Edu NSW ACT Foundation font is used for display/logo purposes
const eduHand = Edu_NSW_ACT_Foundation({
    subsets: ["latin"],
    weight: ["400"],
    variable: "--font-logo", // Custom CSS variable for logo only
    display: "swap",
});

export const metadata: Metadata = {
    title: "Kylon:",
    description:
        "Kylon is a robust backend application designed for distributed systems management, offering streamlined operations and enhanced control.",
    keywords: [
        "Kylon",
        "Distributed Systems",
        "System Management",
        "Backend Application",
        "Cloud Native",
    ],
    openGraph: {
        title: "Kylon: Distributed Systems Management Simplified.",
        description:
            "Kylon is a robust backend application designed for distributed systems management, offering streamlined operations and enhanced control.",
        url: "https://www.yourkylonapp.com", // Placeholder URL, update as needed
        siteName: "Kylon Distributed Systems",
        images: [
            {
                url: "/images/kylon-og-image.png", // Ensure this image path exists
                width: 1200,
                height: 630,
                alt: "Kylon Distributed Systems Management Application",
            },
        ],
    },
    twitter: {
        card: "summary_large_image",
        title: "Kylon: Distributed Systems Management Simplified.",
        description:
            "Kylon is a robust backend application designed for distributed systems management, offering streamlined operations and enhanced control.",
        creator: "@yourcompany", // Update with your Twitter handle
        images: ["/images/kylon-og-image.png"], // Ensure this image path exists
    },
};

export default function RootLayout({
    children,
}: {
    children: React.ReactNode;
}) {
    return (
        <html lang="en" className={`${manrope.variable} ${eduHand.variable}`}>
            <body className="font-manrope">
                <div className="main-content-wrapper">
                    <ToastProvider>{children}</ToastProvider>
                </div>
            </body>
        </html>
    );
}
