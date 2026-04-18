import { List, Paragraph, Section } from "../components/legal-content";
import LegalLayout from "../layout/legal-layout";

export default function TermsPage() {
  return (
    <LegalLayout title="Terms of Service" lastUpdated="April 18, 2026">
      <Paragraph>
        Welcome to <b>ProjectMate</b> (“we”, “our”, “us”). By accessing or using
        our website (<b>https://itsdeployedbyme.dpdns.org</b>), you agree to
        these Terms of Service.
      </Paragraph>
      <Section title="1. Use of the Service">
        <Paragraph>
          ProjectMate is a project management application that allows users to
          create, manage, and track projects and tasks.
        </Paragraph>

        <Paragraph>You agree to:</Paragraph>
        <List
          items={[
            "Use the service only for lawful purposes",
            "Not misuse or attempt to disrupt the service",
            "Provide accurate account information",
          ]}
        />
      </Section>
      <Section title="2. Accounts">
        <Paragraph>
          To use certain features, you must create an account.
        </Paragraph>
        <Paragraph>You are responsible for:</Paragraph>
        <List
          items={[
            "Maintaining the confidentiality of your account credentials",
            "All activities under your account",
          ]}
        />
        <Paragraph>
          We reserve the right to suspend or terminate accounts that violate
          these terms.
        </Paragraph>
      </Section>
      <Section title="3. User Data">
        <Paragraph>
          You retain ownership of any data you create or upload. By using
          ProjectMate, you grant us permission to store and process your data to
          provide the service.
        </Paragraph>
      </Section>
      <Section title="4. Email Communications">
        <Paragraph>By signing up, you agree to receive:</Paragraph>
        <List
          items={[
            "Account-related emails (verification, password reset)",
            "Service notifications (if applicable)",
          ]}
        />
        <Paragraph>
          You can opt out of non-essential emails where possible.
        </Paragraph>
      </Section>
      <Section title="5. Service Availability">
        <Paragraph>
          We do not guarantee that the service will always be available, secure,
          or error-free.
        </Paragraph>
      </Section>
      <Section title="6. Limitation of Liability">
        <Paragraph>
          ProjectMate is provided “as is” without warranties of any kind. We are
          not liable for any damages arising from the use of the service.
        </Paragraph>
      </Section>
      <Section title="7. Termination">
        <Paragraph>
          We may suspend or terminate your access at any time if you violate
          these terms.
        </Paragraph>
      </Section>
      <Section title="8. Changes to Terms">
        <Paragraph>
          We may update these Terms occasionally. Continued use of the service
          means you accept the updated terms.
        </Paragraph>
      </Section>
      <Section title="9. Contact">
        <Paragraph>
          For any questions, contact us at:{" "}
          <a href="mailto:hello@contact.itsdeployedbyme.dpdns.org">
            hello@contact.itsdeployedbyme.dpdns.org
          </a>
        </Paragraph>
      </Section>
    </LegalLayout>
  );
}
